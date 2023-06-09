package index

import (
	"bytes"
	"github.com/google/btree"
	"kv-go/data"
)

// Indexer 抽象索引接口，后续会引用其他数据结构
type Indexer interface {
	//Put 向索引中存储key，对应的数据位置信息
	Put(key []byte, pos *data.LogRecordPos) bool
	//Get 根据key取出对应的索引位置信息
	Get(key []byte) *data.LogRecordPos
	//Delete 根据key删除对应的索引位置信息
	Delete(key []byte) bool
}

type IndexType = int8

const (
	//Btree BTree 索引
	Btree IndexType = iota + 1

	//ART 自适应基数树索引
	ART
)

// NewIndexer 根据类型初始化索引
func NewIndexer(typ IndexType) Indexer {
	switch typ {
	case Btree:
		return NewIndexer(typ)
	case ART:
		//todo
		return nil
	default:
		panic("unsupported index type")
	}
}

type Item struct {
	key []byte
	pos *data.LogRecordPos
}

func (ai *Item) Less(bi btree.Item) bool {
	return bytes.Compare(ai.key, bi.(*Item).key) == -1
}
