package fio

import "os"

// FileIO 标准系统文件 IO
type FileIO struct {
	fd *os.File
}

// NewFileIOManager 初始化标准文件
func NewFileIOManager(fileName string) (*FileIO, error) {
	fd, err := os.OpenFile(
		fileName,
		os.O_CREATE|os.O_RDWR|os.O_APPEND,
		DataFilePerm,
	)
	if err != nil {
		return nil, err
	}
	return &FileIO{fd: fd}, nil

}

// Read 从文件的给定位置读取对应的数据
func (f *FileIO) Read(b []byte, offset int64) (int, error) {
	return f.fd.ReadAt(b, offset)
}

// Write 写入字节数组到文件中
func (f *FileIO) Write(b []byte) (int, error) {
	return f.fd.Write(b)
}

// Sync 持久化数据
func (f *FileIO) Sync() error {
	return f.fd.Sync()
}

// Close 关闭文件
func (f *FileIO) Close() error {
	return f.fd.Close()
}
