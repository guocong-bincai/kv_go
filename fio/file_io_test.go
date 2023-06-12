package fio

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestNewFileManager(t *testing.T) {
	//在电脑的tmp文件夹下创建一个a.data文件
	fio, err := NewFileIOManager(filepath.Join("/tmp", "a.data"))
	assert.Nil(t, err)
	assert.NotNil(t, fio)
}

// destroyFile 要及时删除文件
func destroyFile(name string) {
	if err := os.RemoveAll(name); err != nil {
		panic(err)
	}
}

func TestFileIO_Write(t *testing.T) {
	//在电脑的tmp文件夹下创建一个a.data文件
	path := filepath.Join("/tmp", "a.data")
	fio, err := NewFileIOManager(path)
	defer destroyFile(path)
	assert.Nil(t, err)
	assert.NotNil(t, fio)
	//在电脑的tmp文件夹下创建一个a.data文件，写入“”
	n, err := fio.Write([]byte(""))
	assert.Equal(t, 0, n)
	assert.Nil(t, err)
	//在电脑的tmp文件夹下创建一个a.data文件，写入“bitcask kv”
	n, err = fio.Write([]byte("bitcask kv"))
	t.Log(n, err)
	//在电脑的tmp文件夹下创建一个a.data文件，写入“storage”
	n, err = fio.Write([]byte("storage"))
	t.Log(n, err)
}

func TestFileIO_Read(t *testing.T) {
	//在电脑的tmp文件夹下创建一个001.data文件
	path := filepath.Join("/tmp", "a.data")
	fio, err := NewFileIOManager(path)
	defer destroyFile(path)
	assert.Nil(t, err)
	assert.NotNil(t, fio)
	//在电脑的tmp文件夹下创建一个a.data文件，写入“key-a”
	_, err = fio.Write([]byte("key-a"))
	assert.Nil(t, err)
	//在电脑的tmp文件夹下创建一个a.data文件，写入“key-as”
	_, err = fio.Write([]byte("key-as"))
	assert.Nil(t, err)

	//读数据的时候必须要制定一个字节数组
	b := make([]byte, 5)
	//从这个文件第0个位置开始读，读5个字节的数据，n返回的是读取的字节数
	n, err := fio.Read(b, 0)
	t.Log(string(b), n)
	assert.Equal(t, 5, n)
	assert.Equal(t, []byte("key-a"), b)

	//读数据的时候必须要制定一个字节数组
	b1 := make([]byte, 6)
	//从这个文件第5个位置开始读，读6个字节的数据，n返回的是读取的字节数
	n1, err := fio.Read(b1, 5)
	t.Log(string(b1), n1)
	assert.Equal(t, 6, n1)
	assert.Equal(t, []byte("key-as"), b1)
	//01:01:50
}

func TestFileIO_Sync(t *testing.T) {
	//在电脑的tmp文件夹下创建一个001.data文件
	path := filepath.Join("/tmp", "001.data")
	fio, err := NewFileIOManager(path)
	defer destroyFile(path)
	assert.Nil(t, err)
	assert.NotNil(t, fio)

	err = fio.Sync()
	assert.Nil(t, err)
}

func TestFileIO_Close(t *testing.T) {
	//在电脑的tmp文件夹下创建一个001.data文件
	path := filepath.Join("/tmp", "001.data")
	fio, err := NewFileIOManager(path)
	defer destroyFile(path)
	assert.Nil(t, err)
	assert.NotNil(t, fio)

	err = fio.Close()
	assert.Nil(t, err)
}

//在Go语言中，可以使用go test命令来运行测试文件。如果需要全量跑测试文件，可以在命令行中进入测试文件所在的目录，然后执行以下命令：go test -v
//全量跑测试的指令：go test -v ./
