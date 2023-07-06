# bitcast kv_go 数据库

### 打开tmp文件指令
    open /tmp

### 进入tmp文件指令
    cd /tmp

### 浏览tmp文件指令
    ls /tmp

### 项目写入数据的位置
#### /var/folders/lt/7b_qzcbs30s5lx44fst10gk00000gn/T

### 数据格式设计思路
##### 1.首先是使用EncodeLogRecord编码方法，其功能是将LogRecord结构体转换为符合我们日志记录格式的字节数组
##### 2.我们现将header部分的几个字段写入到对应的字节数组中，header的这几个字段的占据的空间如下：
    2.1 crc是uint类型的，占4个字节
    2.2 Type定义为byte类型，只需要1个字节
    2.3 keySize和valueSize是变长的，每一个的最大值是5

### 最后观看时间
