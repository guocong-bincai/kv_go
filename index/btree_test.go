package index

import (
	"fmt"
	"github.com/stretchr/testify/assert"

	"kv-go/data"
	"testing"
)

func TestBTree_Put(t *testing.T) {
	bt := NewBTree()

	res := bt.Put(nil, &data.LogRecordPos{Fid: 1, Offset: 100})
	fmt.Println(res)

	res1 := bt.Put([]byte("a"), &data.LogRecordPos{Fid: 1, Offset: 2})
	fmt.Println(res1)

	res2 := bt.Put([]byte("as"), &data.LogRecordPos{Fid: 1, Offset: 3})
	fmt.Println(res2)
}

func TestBTree_Get(t *testing.T) {
	bt := NewBTree()

	res := bt.Put(nil, &data.LogRecordPos{Fid: 1, Offset: 100})
	fmt.Println(res)

	pos1 := bt.Get(nil)
	assert.Equal(t, uint32(1), pos1.Fid)

	res2 := bt.Put([]byte("a1"), &data.LogRecordPos{Fid: 1, Offset: 3})
	fmt.Println(res2)

	res3 := bt.Put([]byte("a2"), &data.LogRecordPos{Fid: 1, Offset: 3})
	fmt.Println(res3)

}

func TestBTree_Delete(t *testing.T) {
	bt := NewBTree()
	res := bt.Put(nil, &data.LogRecordPos{Fid: 1, Offset: 100})
	fmt.Println(res)

	res3 := bt.Put([]byte("aaa"), &data.LogRecordPos{Fid: 12, Offset: 22})
	fmt.Println(res3)

	res4 := bt.Delete([]byte("aaa"))
	fmt.Println(res4)

}
