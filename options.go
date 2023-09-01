package kv_go

import "os"

type Options struct {
	//数据库数据目录
	DirPath string
	//数据文件的大小
	DataFileSize int64
	//每次写数据是否持久化
	SyncWrites bool
	//索引类型
	IndexerType IndexerType
	// 索引类型
	IndexType IndexerType
}

// IteratorOptions 索引迭代器配置项
type IteratorOptions struct {
	//遍历前缀为制定值的 Key，默认为空
	Prefix []byte
	//是否反响遍历，默认false是正向
	Reverse bool
}

type IndexerType = int8

const (
	//BTree 索引
	BTree IndexerType = iota + 1
	//ART Adpative Radix Tree 自适应基数树索引
	ART
)

var DefaultOptions = Options{
	DirPath:      os.TempDir(),
	DataFileSize: 256 * 1024 * 1024, //256M
	SyncWrites:   false,
	IndexerType:  BTree,
}

var DefaultIteratorOptions = IteratorOptions{
	Prefix:  nil,
	Reverse: false,
}
