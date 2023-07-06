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
