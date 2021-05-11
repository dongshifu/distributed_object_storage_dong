package locate

import (
	"object_storage_dong/lib/rabbitmq"
	"os"
	"strconv"
	"time"
)

func Locate(name string) string {
	//创建一个消息队列
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	// 向dataServers exchange群发对象名字的定位信息
	q.Publish("dataServers", name)
	c := q.Consume()
	go func() { //1s后关闭临时消息队列
		//设置超时机制，避免无止境的等待。
		//1s后没有任何反馈，消息队列关闭，收到一个长度为0的消息，返回一个空字符串
		time.Sleep(time.Second)
		q.Close()
	}()
	//阻塞等待数据服务节点向自己发送反馈消息
	//若在1s内有来自数据服务节点的消息，返回该消息的正文内容，也就是该数据服务节点的监听地址
	msg := <-c
	s, _ := strconv.Unquote(string(msg.Body))
	return s
}

// 检查Locate结果是否为空字符串来判定对象是否存在
func Exist(name string) bool {
	return Locate(name) != ""
}
