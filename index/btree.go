package index

import (
	"bytes"
	"github.com/google/btree"
	"kv-go/data"
	"sort"
	"sync"
)

// BTree 索引，主要封装了go BTree库
// github地址：https://github.com/google/btree
type BTree struct {
	tree *btree.BTree
	lock *sync.RWMutex //btree并发写是不安全的，需要加锁
}

// NewBTree 新建 BTree 索引结构
func NewBTree() *BTree {
	return &BTree{
		tree: btree.New(32),
		lock: new(sync.RWMutex),
	}
}

func (bt *BTree) Put(key []byte, pos *data.LogRecordPos) *data.LogRecordPos {
	it := &Item{key: key, pos: pos}
	bt.lock.Lock()
	oldItem := bt.tree.ReplaceOrInsert(it)
	bt.lock.Unlock()
	if oldItem == nil {
		return nil
	}
	return oldItem.(*Item).pos
}

func (bt *BTree) Get(key []byte) *data.LogRecordPos {
	it := &Item{key: key}
	btreeItem := bt.tree.Get(it)
	if btreeItem == nil {
		return nil
	}
	return btreeItem.(*Item).pos
}

func (bt *BTree) Delete(key []byte) (*data.LogRecordPos, bool) {
	it := &Item{key: key}
	bt.lock.Lock()
	oldItem := bt.tree.Delete(it)
	bt.lock.Unlock()
	if oldItem == nil {
		return nil, false
	}
	return oldItem.(*Item).pos, true
}

func (bt *BTree) Iterator(reverse bool) Iterator {
	if bt.tree == nil {
		return nil
	}
	bt.lock.RLock()
	defer bt.lock.RUnlock()
	return newBTreeIterator(bt.tree, reverse)
}

// BTree 索引迭代器
type btreeIterator struct {
	currIndex int     //当前遍历的下标位置
	reverse   bool    //是否是反向遍历
	values    []*Item //key+位置索引信息
}

func newBTreeIterator(tree *btree.BTree, reverse bool) *btreeIterator {
	var idx int
	values := make([]*Item, tree.Len())

	//将所有的数据存储到数组中
	saveValues := func(it btree.Item) bool {
		values[idx] = it.(*Item)
		idx++
		return true
	}
	//是否需要反向操作
	if reverse {
		tree.Descend(saveValues)
	} else {
		tree.Ascend(saveValues)
	}

	return &btreeIterator{
		currIndex: 0,
		reverse:   reverse,
		values:    values,
	}
}

func (bti *btreeIterator) Rewind() {
	bti.currIndex = 0
}

// Seek 根据穿日的key 查找到第一个大于（或小于）等于的目标key，根据从这个key开始遍历
func (bti *btreeIterator) Seek(key []byte) {
	//还是根据reverse去判断顺序，正序或者反序
	if bti.reverse {
		bti.currIndex = sort.Search(len(bti.values), func(i int) bool {
			return bytes.Compare(bti.values[i].key, key) <= 0
		})
	} else {
		bti.currIndex = sort.Search(len(bti.values), func(i int) bool {
			return bytes.Compare(bti.values[i].key, key) >= 0
		})
	}
}

// Next 跳转到下一个key
func (bti *btreeIterator) Next() {
	bti.currIndex += 1
}

// Valid 是否有效，即是否已经遍历完所有的key，用于退出遍历
func (bti *btreeIterator) Valid() bool {
	return bti.currIndex < len(bti.values)
}

// Key 当前遍历位置的key 数据
func (bti *btreeIterator) Key() []byte {
	//返回当前位置对应的数据
	return bti.values[bti.currIndex].key
}

// Value 当前遍历位置的Value 数据
func (bti *btreeIterator) Value() *data.LogRecordPos {
	return bti.values[bti.currIndex].pos
}

// Close 关闭迭代器，释放相应资源
func (bti *btreeIterator) Close() {
	bti.values = nil
}
