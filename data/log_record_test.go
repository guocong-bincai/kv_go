package data

import (
	"fmt"
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
	rec2 := &LogRecord{
		Key:  []byte("name"),
		Type: LogRecordNormal,
	}
	res2, n2 := EncodeLogRecord(rec2)
	assert.NotNil(t, res2)
	assert.Greater(t, n2, int64(5))
	t.Log(res2)
	t.Log(n2)
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

	headerBuf2 := []byte{19, 93, 63, 113, 0, 8, 0}
	h2, size2 := decodeLogRecordHeader(headerBuf2)
	t.Log(h2)
	//{2511199366 0 4 10} :crc+类型+key长度+value长度
	t.Log(size2)
	assert.Equal(t, uint32(1899978003), h2.crc)
	assert.Equal(t, LogRecordNormal, h2.recordType)
	assert.Equal(t, uint32(4), h2.keySize)
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
	assert.Equal(t, uint32(2532332136), crc)

	// value 为空的情况
	rec2 := &LogRecord{
		Key:  []byte("name"),
		Type: LogRecordNormal,
	}
	headerBuf2 := []byte{19, 93, 63, 113, 0, 8, 0}
	crc2 := getLogRecordCRC(rec2, headerBuf2[crc32.Size:])
	fmt.Println(crc2)
	assert.Equal(t, uint32(1899978003), crc2)
}
