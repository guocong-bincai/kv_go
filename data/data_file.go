package data

import (
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"kv-go/fio"
	"path/filepath"
)

var (
	ErrInvalidCRC = errors.New("invalid crc value,log record may be corrupted")
)

const DataFileNameSuffix = ".data"

// DataFile 数据文件
type DataFile struct {
	FileId    uint32        //文件id
	WriteOff  int64         //文件写到了哪个位置
	IoManager fio.IOManager //io 读写管理
}

// OpenDataFile 打开你的数据文件
func OpenDataFile(dirPath string, filed uint32) (*DataFile, error) {
	//根据文件路径，生成文件名称
	fileName := filepath.Join(dirPath, fmt.Sprintf("%09d", filed)+DataFileNameSuffix)
	//初始化IOManager 管理器接口
	ioManager, err := fio.NewIOManager(fileName)
	if err != nil {
		return nil, err
	}
	return &DataFile{
		FileId:    filed,
		WriteOff:  0,
		IoManager: ioManager,
	}, nil
}

// ReadLogRecord 根据 offset 从数据文件中读取 LogRecord
func (df *DataFile) ReadLogRecord(offset int64) (*LogRecord, int64, error) {
	fileSize, err := df.IoManager.Size()
	if err != nil {
		return nil, 0, err
	}

	//如果读取的最大 header 长度已经超过了文件的长度，则只需要读取到文件的末尾即可
	var headerBytes int64 = maxLogRecordHeaderSize
	if offset+maxLogRecordHeaderSize > fileSize {
		headerBytes = fileSize - offset
	}
	//读取 Header 信息
	headerBuf, err := df.readNBytes(headerBytes, offset)
	if err != nil {
		return nil, 0, err
	}
	//然后对headerBuf进行解码
	header, headSize := decodeLogRecordHeader(headerBuf)
	//下面的两个条件表示读取到了文件末尾，直接返回EOF错误
	if header == nil {
		return nil, 0, io.EOF
	}
	if header.crc == 0 && header.keySize == 0 && header.valueSize == 0 {

	}

	//取出对应的key 和value 的长度
	keySize, valueSize := int64(header.keySize), int64(header.valueSize)
	var recordSize = headSize + keySize + valueSize

	logRecord := &LogRecord{Type: header.recordType}
	//开始读取用户实际存储的 key/value数据
	if keySize > 0 || valueSize > 0 {
		kvBuf, err := df.readNBytes(keySize+valueSize, offset+headSize)
		if err != nil {
			return nil, 0, err
		}
		//解出 key 和value
		logRecord.Key = kvBuf[:keySize]
		logRecord.Value = kvBuf[keySize:]
	}

	//校验数据的有效性
	crc := getLogRecordCRC(logRecord, headerBuf[crc32.Size:headSize])
	if crc != header.crc {
		return nil, 0, ErrInvalidCRC
	}
	//上述检测都通过了才说明正确数值然后返回
	return logRecord, recordSize, nil
}

func (df *DataFile) Write(buf []byte) error {
	n, err := df.IoManager.Write(buf)
	if err != nil {
		return err
	}
	df.WriteOff += int64(n)
	return nil
}

func (df *DataFile) Sync() error {
	return df.IoManager.Sync()
}

func (df *DataFile) Close() error {
	return df.IoManager.Close()
}

// 制定读取n个字节，然后调用Read 返回一个数组
func (df *DataFile) readNBytes(n int64, offset int64) (b []byte, err error) {
	b = make([]byte, n)
	_, err = df.IoManager.Read(b, offset)
	return
}
