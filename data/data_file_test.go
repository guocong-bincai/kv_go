package data

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestDataFile(t *testing.T) {
	datafile, err := OpenDataFile(os.TempDir(), 0)
	assert.Nil(t, err)
	assert.NotNil(t, datafile)
}

func TestDataFile_Write(t *testing.T) {
	datafile, err := OpenDataFile(os.TempDir(), 1)
	assert.Nil(t, err)
	assert.NotNil(t, datafile)

	err = datafile.Write([]byte("aaa"))
	assert.Nil(t, err)
}

func TestDataFile_Close(t *testing.T) {
	datafile, err := OpenDataFile(os.TempDir(), 3)
	assert.Nil(t, err)
	assert.NotNil(t, datafile)

	err = datafile.Write([]byte("aaa"))
	assert.Nil(t, err)

	err = datafile.Close()
	assert.Nil(t, err)
}

func TestDataFile_Sync(t *testing.T) {
	datafile, err := OpenDataFile(os.TempDir(), 4)
	assert.Nil(t, err)
	assert.NotNil(t, datafile)

	err = datafile.Write([]byte("aaa"))
	assert.Nil(t, err)

	err = datafile.Sync()
	assert.Nil(t, err)
}

func TestDataFile_ReadLogRecord(t *testing.T) {
	dataFile, err := OpenDataFile(os.TempDir(), 222)
	assert.Nil(t, err)
	assert.NotNil(t, dataFile)

	//只有一条LogRecord
	rec1 := &LogRecord{
		Key:   []byte("name"),
		Value: []byte("bitcask kv go"),
	}
	res1, size1 := EncodeLogRecord(rec1)
	err = dataFile.Write(res1)
	assert.Nil(t, err)

	readRec1, readSize1, err := dataFile.ReadLogRecord(0)
	assert.Nil(t, err)
	assert.Equal(t, rec1, readRec1)
	assert.Equal(t, size1, readSize1)

	//多条LogEecord，从不同的位置读取
	rec2 := &LogRecord{
		Key:   []byte("name"),
		Value: []byte("a new value"),
	}
	res2, size2 := EncodeLogRecord(rec2)
	err = dataFile.Write(res2)
	assert.Nil(t, err)
	t.Log(size2)
}
