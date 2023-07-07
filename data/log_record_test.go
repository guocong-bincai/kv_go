package data

import (
	"github.com/stretchr/testify/assert"
	"hash/crc32"
	"testing"
)

func TestEncodeLogRecord(t *testing.T) {
	// 正常情况
	rec1 := &LogRecord{
		Key:   []byte("name"),
		Value: []byte("bitcask-go"),
		Type:  LogRecordNormal,
	}
	res1, n1 := EncodeLogRecord(rec1)
	assert.NotNil(t, res1)
	//crc + type 固定的5个字节，所以结果一定是大于5的
	assert.Greater(t, n1, int64(5))
	t.Log(res1)
	t.Log(n1)
	// value 为空的情况

	// 对Deleted 情况的测试
}

func TestDecodeLogRecordHeader(t *testing.T) {
	headerBuf1 := []byte{134, 220, 173, 149, 0, 8, 20}
	h1, size1 := decodeLogRecordHeader(headerBuf1)
	t.Log(h1)
	//{2511199366 0 4 10} :crc+类型+key长度+value长度
	t.Log(size1)
	assert.Equal(t, uint32(2511199366), h1.crc)
	assert.Equal(t, LogRecordNormal, h1.recordType)
	assert.Equal(t, uint32(4), h1.keySize)
	assert.Equal(t, uint32(10), h1.valueSize)
}

func TestGetLogRecordCRC(t *testing.T) {
	// 正常情况
	rec1 := &LogRecord{
		Key:   []byte("name"),
		Value: []byte("bitcask-go"),
		Type:  LogRecordNormal,
	}
	headerBuf1 := []byte{134, 220, 173, 149, 0, 8, 20}
	crc := getLogRecordCRC(rec1, headerBuf1[crc32.Size:])
	t.Log(crc)
}
