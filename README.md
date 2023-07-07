# bitcast kv_go 数据库

### 打开tmp文件指令
    open /tmp

### 进入tmp文件指令
    cd /tmp

### 浏览tmp文件指令
    ls /tmp

### 项目写入数据的位置
#### /var/folders/lt/7b_qzcbs30s5lx44fst10gk00000gn/T

### LogRecord编码-数据格式设计思路
##### 1.首先是使用EncodeLogRecord编码方法，其功能是将LogRecord结构体转换为符合我们日志记录格式的字节数组
##### 2.我们现将header部分的几个字段写入到对应的字节数组中，header的这几个字段的占据的空间如下：
    2.1 crc是uint类型的，占4个字节
    2.2 Type定义为byte类型，只需要1个字节
    2.3 keySize和valueSize是变长的，每一个的最大值是5

### 关于CRC校验的信息
    CRC（Cyclic Redundancy Check）校验是一种常用的错误检测技术，用于验证数据在传输过程中是否发生了错误或损坏。
    CRC 校验通过计算数据的循环冗余校验码（CRC 值），将其附加到数据中，接收方可以使用相同的算法计算 CRC 值并与接收到的 CRC 值进行比较，以确定数据的完整性。
    CRC 校验的基本原理是将数据看作二进制位序列，并使用一个预定义的生成多项式进行除法运算。
    生成多项式通常是一个固定的二进制数，如 CRC-16、CRC-32 等。
    通过除法运算，得到的余数就是 CRC 值，将其附加到原始数据中进行传输。

### LogRecord解码-数据格式设计思路
    从数据文件中读取日志记录LogRecord时，首先会按照固定大小读取header部分的字节数，
    然后对其进行解码，主要是根据编码时的对应长度获取CRC校验值，Type，key size，value size，
    要是根据编码时的对应长度CRC校验值，type，key size，value size。
    然后再根据key size 和value size去除实际的key/value数据。
    最后需要校验读取出的crc值是否和LogRecord对应的crc值是否相等，
    如果不想等的话则说明这条数据存在村务，那么需要返回对应的错误信息。

### 最后观看时间
09 00:33:00