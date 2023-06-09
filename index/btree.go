package index

import (
	"github.com/google/btree"
	"kv-go/data"
	"sync"
)

// BTree 索引，主要封装了go BTree库
// github地址：https://github.com/google/btree
type BTree struct {
	tree *btree.BTree
	lock *sync.RWMutex //btree并发写是不安全的，需要加锁
}

// NewBTree 初始化BTree索引结构
func NewBTree() *BTree {
	return &BTree{
		tree: btree.New(32),
		lock: new(sync.RWMutex),
	}
}

// PutV1 存入数据
func (bt *BTree) PutV1(key []byte, pos *data.LogRecordPos) bool {
	it := Item{key: key, pos: pos}
	bt.lock.Lock()
	bt.tree.ReplaceOrInsert(&it)
	bt.lock.Unlock()

	return true
}

// Put 存入数据
func (bt *BTree) Put(key []byte, pos *data.LogRecordPos) *data.LogRecordPos {
	it := Item{key: key, pos: pos}
	bt.lock.Lock()
	oldItem := bt.tree.ReplaceOrInsert(&it)
	bt.lock.Unlock()
	if oldItem == nil {
		return nil
	}
	return oldItem.(*Item).pos
}

// GetV1 查找数据
func (bt *BTree) GetV1(key []byte) *data.LogRecordPos {
	it := Item{key: key}
	btreeItem := bt.tree.Get(&it)
	if btreeItem == nil {
		return nil
	}
	return btreeItem.(*Item).pos

}

// Get 查找数据
func (bt *BTree) Get(key []byte) *data.LogRecordPos {
	it := Item{key: key}
	btreeItem := bt.tree.Get(&it)
	if btreeItem == nil {
		return nil
	}
	return btreeItem.(*Item).pos

}

// DeleteV1 删除数据
func (bt *BTree) DeleteV1(key []byte) bool {
	it := Item{key: key}
	bt.lock.Lock()
	olderItem := bt.tree.Delete(&it)
	bt.lock.Unlock()
	if olderItem == nil {
		return false
	}
	return true
}

// Delete 删除数据
func (bt *BTree) Delete(key []byte) *data.LogRecordPos {
	it := Item{
		key: key,
	}
	bt.lock.Lock()
	olderItem := bt.tree.Delete(&it)
	bt.lock.Unlock()
	if olderItem != nil {
		return nil
	}

	return olderItem.(*Item).pos
}
