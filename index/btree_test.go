package index

import (
	"github.com/stretchr/testify/assert"

	"kv-go/data"
	"testing"
)

func TestBTree_Iterator(t *testing.T) {
	bt1 := NewBTree()
	//1.Btree 为空的情况
	iter1 := bt1.Iterator(false)
	assert.Equal(t, false, iter1.Valid())

	// 2.BTree 有数据的情况下
	bt1.Put([]byte("code"), &data.LogRecordPos{
		Fid:    1,
		Offset: 10,
	})
	iter2 := bt1.Iterator(false)
	assert.Equal(t, true, iter2.Valid())
	t.Log(iter2.Key())
	t.Log(iter2.Value())
	iter2.Next()
	assert.Equal(t, false, iter2.Valid())

	//有多条数据
	bt1.Put([]byte("aaa"), &data.LogRecordPos{Fid: 1, Offset: 10})
	bt1.Put([]byte("bbb"), &data.LogRecordPos{Fid: 1, Offset: 10})
	bt1.Put([]byte("ccc"), &data.LogRecordPos{Fid: 1, Offset: 10})
	bt1.Put([]byte("ddd"), &data.LogRecordPos{Fid: 1, Offset: 10})
	iter3 := bt1.Iterator(false)
	for iter3.Rewind(); iter3.Valid(); iter3.Next() {
		t.Log("key= ", string(iter3.Key()))
	}

	iter4 := bt1.Iterator(false)
	for iter4.Rewind(); iter4.Valid(); iter4.Next() {
		t.Log("key= ", string(iter4.Key()))
	}

	//4.测试 seek
	iter5 := bt1.Iterator(false)
	for iter5.Seek([]byte("cc")); iter5.Valid(); iter5.Next() {
		//t.Log(string(iter5.Key()))
		assert.NotNil(t, iter5.Key())
	}

	//5.反向遍历 seek
	iter6 := bt1.Iterator(true)
	iter6.Seek([]byte("zzz"))
	t.Log(string(iter6.Key()))
}
