package data

// readme:关于数据文件，数据目录的文件

// LogRecordPos 数据内存索引的数据结构，主要是描述数据在磁盘的位置
type LogRecordPos struct {
	Fid    uint32 //文件id，表示的是将数据存储到那个位置中
	Offset int64  //偏移，表示将数据存储到了数据文件中的那个位置
}
