package locate

import (
	"object_storage_dong/lib/rabbitmq"
	"os"
	"strconv"
)

// 实际定位对象
func Locate(name string) bool {
	// os.Stat访问磁盘上对应的文件名
	_, err := os.Stat(name)
	// 判断文件名是否存在，存在返回true,失败返回false
	return !os.IsNotExist(err)
}

// 用于监听定位消息
func StartLocate() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	// 绑定dataServers exchange
	q.Bind("dataServers")
	// 返回一个channel
	c := q.Consume()
	// range遍历channel,接收消息
	for msg := range c {
		// JSON编码使得对象名字上有一对双引号，使用strconv.Unquote将输入的字符串前后的双引号去除并作为结果返回
		object, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		if Locate(os.Getenv("STORAGE_ROOT") + "/objects/" + object) {
			//文件存在，调用Send方法向消息的发送方返回本服务节点的监听地址，表示该对象存在于本服务节点上。
			q.Send(msg.ReplyTo, os.Getenv("LISTEN_ADDRESS"))
		}
	}
}
