package kv_go

import (
	"kv-go/data"
	"kv-go/index"
	"sync"
)

type DB struct {
	options    Options
	mu         *sync.RWMutex
	activeFile *data.DataFile            //当前活跃数据文件，可以用于写入
	olderFile  map[uint32]*data.DataFile //旧的数据文件，只能用于读
	index      index.Indexer             //内存索引
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
	if ok := db.index.Put(key, pos); !ok {
		return ErrIndexUpdateFailed
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

	//根据文件id 找到对应的数据文件
	var dataFile *data.DataFile
	if db.activeFile.FileId == logRecordPos.Fid {
		dataFile = db.activeFile
	} else {
		dataFile = db.olderFile[logRecordPos.Fid]
	}
	//数据文件为空
	if dataFile == nil {
		return nil, ErrDataFileFound
	}
	//根据偏移读量取对应的数据
	logRecord, err := dataFile.ReadLogRecord(logRecordPos.Offset)
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
	dataFile, err := data.OpenDataFile(db.options.DirPath, initialFileId)
	if err != nil {
		return err
	}
	db.activeFile = dataFile
	return nil
}

//28:28
