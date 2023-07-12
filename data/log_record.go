package data

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
)

// readme:关于数据文件，数据目录的文件

type LogRecordType = byte

const (
	LogRecordNormal LogRecordType = iota
	LogRecordDeleted
)

// 头部：crc=4byte,type =1byte ,keySize=变长max=5,valueSize=变长max=5
const maxLogRecordHeaderSize = binary.MaxVarintLen32*2 + 5

// LogRecord 写入到数据文件的记录，之所以叫日志，是因为数据文件中的数据是追加写入的，类似日志的格式
type LogRecord struct {
	Key   []byte
	Value []byte
	Type  LogRecordType
}

// LogRecord 的头部信息
type logRecordHeader struct {
	crc        uint32        //crc 校验值
	recordType LogRecordType //标识 LogRecord的类型
	keySize    uint32        //key的长度
	valueSize  uint32        //value的长度
}

// LogRecordPos 数据内存索引的数据结构，主要是描述数据在磁盘的位置
type LogRecordPos struct {
	Fid    uint32 //文件id，表示的是将数据存储到那个位置中
	Offset int64  //偏移，表示将数据存储到了数据文件中的那个位置
}

// EncodeLogRecord 对 LogRecord 进行编码，返回字节数组及长度
//
//	+-------------+-------------+-------------+--------------+-------------+--------------+
//	| crc 校验值  |  type 类型   |    key size |   value size |      key    |      value   |
//	+-------------+-------------+-------------+--------------+-------------+--------------+
//	    4字节          1字节        变长（最大5）   变长（最大5）     变长           变长

func EncodeLogRecord(logRecord *LogRecord) ([]byte, int64) {
	//初始化一个header 部分的字节数组
	header := make([]byte, maxLogRecordHeaderSize)

	//第5个字节存储，Type
	header[4] = logRecord.Type
	var index = 5
	// 5字节之后，存储的是key和value 的长度信息
	// 使用变长类型，节省空间
	// 从第5个字节开始写入key的大小，注意变量递增
	index += binary.PutVarint(header[index:], int64(len(logRecord.Key)))
	index += binary.PutVarint(header[index:], int64(len(logRecord.Value)))

	var size = index + len(logRecord.Key) + len(logRecord.Value)
	encBytes := make([]byte, size)

	//将header 部分的内容拷贝过来
	copy(encBytes[:index], header[:index])
	//将key 和value 数据拷贝到字节数组中
	copy(encBytes[:index], logRecord.Key)
	copy(encBytes[index+len(logRecord.Key):], logRecord.Value)

	//对整个LogRecord 的数据进行crc校验
	crc := crc32.ChecksumIEEE(encBytes[4:])
	binary.LittleEndian.PutUint32(encBytes[:4], crc)
	return encBytes, int64(size)
}

// 对字节数组中的 Header 信息进行解码
func decodeLogRecordHeader(buf []byte) (*logRecordHeader, int64) {
	if len(buf) <= 4 {
		return nil, 0
	}

	header := &logRecordHeader{
		crc:        binary.LittleEndian.Uint32(buf[:4]),
		recordType: buf[4],
	}

	var index = 5
	//取出实际的key size
	keySize, n := binary.Varint(buf[index:])
	header.keySize = uint32(keySize)
	index += n

	//取出实际的 value size
	valueSize, n := binary.Varint(buf[index:])
	header.valueSize = uint32(valueSize)
	index += n
	fmt.Printf("crc: %d", header.crc)
	return header, int64(index)
}

func getLogRecordCRC(lr *LogRecord, header []byte) uint32 {
	if lr == nil {
		return 0
	}
	crc := crc32.ChecksumIEEE(header[:])
	crc = crc32.Update(crc, crc32.IEEETable, lr.Key)
	crc = crc32.Update(crc, crc32.IEEETable, lr.Value)
	return crc
}
