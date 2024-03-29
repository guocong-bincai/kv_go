package kv_go

import (
	"bytes"
	"kv-go/index"
)

// Iterator 迭代器
type Iterator struct {
	indexIter index.Iterator
	db        *DB
	options   IteratorOptions
}

func (db *DB) NewIterator(opts IteratorOptions) *Iterator {
	indexIter := db.index.Iterator(opts.Reverse)
	return &Iterator{
		db:        db,
		indexIter: indexIter,
		options:   opts,
	}
}

// Rewind 重新回到迭代器的起点，即第一个数据
func (it *Iterator) Rewind() {
	it.indexIter.Rewind()
}

// Seek 根据传入的key 查找到第一个大于（或小于）等于的目标key
func (it *Iterator) Seek(key []byte) {
	it.indexIter.Seek(key)
	//加上跳转的逻辑
	it.skipToNext()
}

// Next 跳转到下一个key
func (it *Iterator) Next() {
	it.indexIter.Next()
	//加上跳转的逻辑
	it.skipToNext()
}

// Valid 是否有效，即是否已经遍历完了所有的key 用于退出遍历
func (it *Iterator) Valid() bool {
	return it.indexIter.Valid()
}

// Key 当前遍历位置的key数据
func (it *Iterator) Key() []byte {
	return it.indexIter.Key()
}

// Value 当前遍历位置的Value数据
func (it *Iterator) Value() ([]byte, error) {
	//拿到位置的索引信息
	logRecordsPos := it.indexIter.Value()
	//根据位置的索引信息拿到value
	it.db.mu.RLock()
	defer it.db.mu.RUnlock()
	return it.db.getValueByPosition(logRecordsPos)
}

// Close 关闭迭代器，释放相应资源
func (it *Iterator) Close() {
	it.indexIter.Close()
}

func (it *Iterator) skipToNext() {
	prefixLen := len(it.options.Prefix)
	if prefixLen == 0 {
		return
	}

	for ; it.indexIter.Valid(); it.indexIter.Next() {
		key := it.indexIter.Key()
		if prefixLen <= len(key) && bytes.Compare(it.options.Prefix, key[:prefixLen]) == 0 {
			break
		}
	}
}
