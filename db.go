package kv_go

import (
	"errors"
	"fmt"
	"github.com/gofrs/flock"
	"io"
	"kv-go/data"
	"kv-go/fio"
	"kv-go/index"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const (
	seqNoKey     = "seq.no"
	fileLockName = "flock"
)

type DB struct {
	options     Options
	mu          *sync.RWMutex
	fileIds     []int                     //文件id，只能在加载文件索引的时候使用，不能在其他的地方更新和使用
	activeFile  *data.DataFile            //当前活跃数据文件，可以用于写入
	olderFile   map[uint32]*data.DataFile //旧的数据文件，只能用于读
	index       index.Indexer             //内存索引
	reclaimSize int64                     // 表示有多少数据是无效的
	fileLock    *flock.Flock              // 文件锁保证多进程之间的互斥
	seqNo       uint64                    // 事务序列号，全局递增
	olderFiles  map[uint32]*data.DataFile // 旧的数据文件，只能用于读
}

// Open 打开 bitcask 存储引擎实例
func Open(options Options) (*DB, error) {
	//03:35
	//对用户传入的配置项进行校验
	if err := checkOptions(options); err != nil {
		return nil, err
	}
	//判断数据目录是否存在，如果不存在，则创建这个目录
	if _, err := os.Stat(options.DirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(options.DirPath, os.ModePerm); err != nil {
			return nil, err
		}
	}
	//初始化 DB实例结构体
	db := &DB{
		options:   options,
		mu:        new(sync.RWMutex),
		olderFile: make(map[uint32]*data.DataFile),
		index:     index.NewIndexer(options.IndexType, options.DirPath, options.SyncWrites),
	}
	//加载数据文件
	if err := db.loadDataFiles(); err != nil {
		return nil, err
	}

	//从数据文件中加载索引
	if err := db.loadIndexFromDataFile(); err != nil {
		return nil, err
	}
	return db, nil
}

// Close 关闭数据库
func (db *DB) Close() error {
	defer func() {
		// 释放文件锁
		if err := db.fileLock.Unlock(); err != nil {
			panic(fmt.Sprintf("failed to unlock the directory, %v", err))
		}
		// 关闭索引
		if err := db.index.Close(); err != nil {
			panic(fmt.Sprintf("failed to close index"))
		}
	}()
	if db.activeFile == nil {
		return nil
	}
	db.mu.Lock()
	defer db.mu.Unlock()

	// 保存当前事务序列号
	seqNoFile, err := data.OpenSeqNoFile(db.options.DirPath)
	if err != nil {
		return err
	}
	record := &data.LogRecord{
		Key:   []byte(seqNoKey),
		Value: []byte(strconv.FormatUint(db.seqNo, 10)),
	}
	encRecord, _ := data.EncodeLogRecord(record)
	if err := seqNoFile.Write(encRecord); err != nil {
		return err
	}
	if err := seqNoFile.Sync(); err != nil {
		return err
	}

	//	关闭当前活跃文件
	if err := db.activeFile.Close(); err != nil {
		return err
	}
	// 关闭旧的数据文件
	for _, file := range db.olderFiles {
		if err := file.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Put 写入 key/value 数据，key不能为空
func (db *DB) Put(key []byte, value []byte) error {
	//判断key是否有效
	if len(key) == 0 {
		return ErrKeyIsEmpty
	}

	//构造 LogRecord 结构体
	logRecord := &data.LogRecord{
		Key:   key,
		Value: value,
		Type:  data.LogRecordNormal,
	}
	pos, err := db.appendLogRecord(logRecord)
	if err != nil {
		return err
	}
	//更新内存索引
	// 更新内存索引
	if oldPos := db.index.Put(key, pos); oldPos != nil {
		db.reclaimSize += int64(oldPos.Size)
	}

	return nil
}

// Get 根据key读取数据
func (db *DB) Get(key []byte) ([]byte, error) {
	//判断key的有效性
	db.mu.RLock()
	defer db.mu.RUnlock()
	if len(key) == 0 {
		return nil, ErrKeyIsEmpty
	}

	//从内存数据结构中取出 key 对应的索引信息
	logRecordPos := db.index.Get(key)
	//如果key不在内存索引中，说明key不存在
	if logRecordPos == nil {
		return nil, ErrKeyNotFound
	}

	// 从数据文件中获取value
	return db.getValueByPosition(logRecordPos)
}

// 根据索引消息获取对应的value
func (db *DB) getValueByPosition(pos *data.LogRecordPos) ([]byte, error) {
	//根据文件id 找到对应的数据文件
	var dataFile *data.DataFile
	if db.activeFile.FileId == pos.Fid {
		dataFile = db.activeFile
	} else {
		dataFile = db.olderFile[pos.Fid]
	}
	//数据文件为空
	if dataFile == nil {
		return nil, ErrDataFileFound
	}
	//根据偏移读量取对应的数据
	logRecord, _, err := dataFile.ReadLogRecord(pos.Offset)
	if err != nil {
		return nil, err
	}

	//判断logRecord的类型，看是否是被删除的
	if logRecord.Type == data.LogRecordDeleted {
		return nil, ErrKeyNotFound
	}

	//返回实际存储的数据
	return logRecord.Value, nil
}

// appendLogRecord 追加写数据到活跃文件中
func (db *DB) appendLogRecord(logRecord *data.LogRecord) (*data.LogRecordPos, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	//判断当前活跃数据文件是否存在，因为数据库在没有写入的时候是没有文件生成的
	//如果为空则初始化数据文件
	if db.activeFile == nil {
		if err := db.setActiveDataFile(); err != nil {
			return nil, err
		}
	}

	//写入数据编码
	encRecord, size := data.EncodeLogRecord(logRecord)
	//如果写入的数据已经到达了活跃文件的阀值，则关闭活跃文件，并打开新的文件
	if db.activeFile.WriteOff+size > db.options.DataFileSize {
		//先持久话数据文件，保证已有的数据持久化到磁盘当中
		if err := db.activeFile.Sync(); err != nil {
			return nil, err
		}
		//当前活跃文件转换为旧的数据文件
		db.olderFile[db.activeFile.FileId] = db.activeFile
		//打开新的数据文件
		if err := db.setActiveDataFile(); err != nil {
			return nil, err
		}
	}
	writeOff := db.activeFile.WriteOff
	if err := db.activeFile.Write(encRecord); err != nil {
		return nil, err
	}

	//根据用户配置决定是否持久化
	if db.options.SyncWrites {
		if err := db.activeFile.Sync(); err != nil {
			return nil, err
		}
	}

	//构造内存索引信息
	pos := &data.LogRecordPos{Fid: db.activeFile.FileId, Offset: writeOff}
	return pos, nil
}

// 设置当前活跃文件
// 在访问此方法前必须持有互斥锁
func (db *DB) setActiveDataFile() error {
	var initialFileId uint32 = 0
	//如果活跃文件不为空，那么当前活跃文件的id是以前的活跃文件id+1
	if db.activeFile != nil {
		initialFileId = db.activeFile.FileId + 1
	}
	//打开新的数据文件
	dataFile, err := data.OpenDataFile(db.options.DirPath, initialFileId, fio.StandardFIO)
	if err != nil {
		return err
	}
	db.activeFile = dataFile
	return nil
}

//28:28
//5集完

func checkOptions(options Options) error {
	if options.DirPath == "" {
		return errors.New("database dir path is empty")
	}
	if options.DataFileSize <= 0 {
		return errors.New("database data file size must be greater than 0")
	}
	return nil
}

// 加载对应的数据文件
func (db *DB) loadDataFiles() error {
	dirEntries, err := os.ReadDir(db.options.DirPath)
	if err != nil {
		return err
	}

	var fileIds []int
	//遍历目录中的所有文件，找到所有以.data结尾的文件
	for _, entry := range dirEntries {
		if strings.HasSuffix(entry.Name(), data.DataFileNameSuffix) {
			splitNames := strings.Split(entry.Name(), ".")
			//数据目录有可能被损坏了
			filed, err := strconv.Atoi(splitNames[0])
			if err != nil {
				return ErrDataDirectoryCorrupted
			}
			//将对应的文件id存放到列表当中
			fileIds = append(fileIds, filed)
		}
	}
	//对文件Id 进行排序，从小到大依次加载
	sort.Ints(fileIds)
	db.fileIds = fileIds

	//遍历每个文件id，打开对应的数据文件
	for i, fid := range fileIds {
		ioType := fio.StandardFIO
		dataFile, err := data.OpenDataFile(db.options.DirPath, uint32(fid), ioType)
		if err != nil {
			return err
		}
		if i == len(fileIds)-1 {
			//最后一个，id是最大的，说明是当前活跃文件
			db.activeFile = dataFile
		} else {
			//说明是旧的数据文件
			db.olderFile[uint32(fid)] = dataFile
		}
	}
	return nil
}

// 从数据文件中加载索引
// 遍历文件中的所有记录，宾更新到数据内存索引中
func (db *DB) loadIndexFromDataFile() error {
	//没有文件，说明数据库是空的，直接返回
	if len(db.fileIds) == 0 {
		return nil
	}

	//遍历所有的文件id，处理文件中的记录
	for i, fid := range db.fileIds {
		var fileId = uint32(fid)
		var dataFile *data.DataFile
		if fileId == db.activeFile.FileId {
			dataFile = db.activeFile
		} else {
			dataFile = db.olderFile[fileId]
		}

		var offset int64 = 0
		//循环处理,读取文件
		for {
			logReocrd, size, err := dataFile.ReadLogRecord(offset)
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}

			//构造内存索引并保存
			logRecordPos := &data.LogRecordPos{Fid: fileId, Offset: offset}
			//如果是删除类型，直接删除掉。
			if logReocrd.Type == data.LogRecordDeleted {
				db.index.Delete(logReocrd.Key)
			} else {
				db.index.Put(logReocrd.Key, logRecordPos)
			}

			//递增 offset，下一次从新的位置开始读取
			offset += size
		}

		//如果是当前活跃文件，更细你这个文件的 WriteOff
		if i == len(db.fileIds)-1 {
			db.activeFile.WriteOff = offset
		}
	}
	return nil
}

// Delete 根据key 删除对应的数据
func (db *DB) Delete(key []byte) error {
	if len(key) == 0 {
		return ErrKeyIsEmpty
	}
	//先检查 key是否存在，如果不存在直接返回
	if pos := db.index.Get(key); pos == nil {
		return nil
	}

	//构造 LogRecord，标识其是被删除的
	logRecord := &data.LogRecord{Key: key, Type: data.LogRecordDeleted}
	//写入到数据文件中
	_, err := db.appendLogRecord(logRecord)
	if err != nil {
		return nil
	}
	//从内存索引中将对应的key删除
	_, ok := db.index.Delete(key)
	if !ok {
		return ErrIndexUpdateFailed
	}
	return nil
}

//第六集 04:38
