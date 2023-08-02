package index

import (
	"github.com/stretchr/testify/assert"

	"kv-go/data"
	"testing"
)

//func TestBTree_Put(t *testing.T) {
//	bt := NewBTree()
//
//	res := bt.Put(nil, &data.LogRecordPos{Fid: 1, Offset: 100})
//	fmt.Println(res)
//
//	res1 := bt.Put([]byte("a"), &data.LogRecordPos{Fid: 1, Offset: 2})
//	fmt.Println(res1)
//
//	res2 := bt.Put([]byte("as"), &data.LogRecordPos{Fid: 1, Offset: 3})
//	fmt.Println(res2)
//}
//
//func TestBTree_PutV1(t *testing.T) {
//	bt := NewBTree()
//
//	res := bt.PutV1(nil, &data.LogRecordPos{Fid: 1, Offset: 100})
//	assert.True(t, res)
//
//}
//
//func TestBTree_Get(t *testing.T) {
//	bt := NewBTree()
//
//	res := bt.Put(nil, &data.LogRecordPos{Fid: 1, Offset: 100})
//	fmt.Println(res)
//
//	pos1 := bt.Get(nil)
//	assert.Equal(t, uint32(1), pos1.Fid)
//
//	res2 := bt.Put([]byte("a1"), &data.LogRecordPos{Fid: 1, Offset: 3})
//	fmt.Println(res2)
//
//	res3 := bt.Put([]byte("a2"), &data.LogRecordPos{Fid: 1, Offset: 3})
//	fmt.Println(res3)
//
//}
//
//func TestBTree_GetV1(t *testing.T) {
//	bt := NewBTree()
//	res1 := bt.PutV1(nil, &data.LogRecordPos{Fid: 1, Offset: 100})
//	assert.True(t, res1)
//
//	pos1 := bt.GetV1(nil)
//	assert.Equal(t, uint32(1), pos1.Fid)
//	assert.Equal(t, int64(100), pos1.Offset)
//}
//
////func TestBTree_Delete(t *testing.T) {
////	bt := NewBTree()
////	res := bt.Put(nil, &data.LogRecordPos{Fid: 1, Offset: 100})
////	fmt.Println(res)
////
////	res3 := bt.Put([]byte("aaa"), &data.LogRecordPos{Fid: 12, Offset: 22})
////	fmt.Println(res3)
////
////	res4 := bt.Delete([]byte("aaa"))
////	fmt.Println(res4)
////
////}
//
//func TestBTree_DeleteV1(t *testing.T) {
//	bt := NewBTree()
//
//	res := bt.PutV1(nil, &data.LogRecordPos{Fid: 1, Offset: 100})
//	assert.True(t, res)
//
//	res4 := bt.DeleteV1(nil)
//	assert.True(t, res4)
//
//	res1 := bt.PutV1([]byte("aaa"), &data.LogRecordPos{Fid: 1, Offset: 100})
//	assert.True(t, res1)
//
//	res2 := bt.DeleteV1([]byte("aaa"))
//	assert.True(t, res2)
//
//}

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
}
