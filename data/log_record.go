package data

// readme:关于数据文件，数据目录的文件

type LogRecordType = byte

const (
	LogRecordNormal LogRecordType = iota
	LogRecordDeleted
)

// LogRecord 写入到数据文件的记录，之所以叫日志，是因为数据文件中的数据是追加写入的，类似日志的格式
type LogRecord struct {
	Key   []byte
	Value []byte
	Type  LogRecordType
}

// LogRecordPos 数据内存索引的数据结构，主要是描述数据在磁盘的位置
type LogRecordPos struct {
	Fid    uint32 //文件id，表示的是将数据存储到那个位置中
	Offset int64  //偏移，表示将数据存储到了数据文件中的那个位置
}

// EncodeLogRecord 对 LogRecord 进行编码，返回字节数组及长度
func EncodeLogRecord(d *LogRecord) ([]byte, int64) {
	return nil, 0
}
