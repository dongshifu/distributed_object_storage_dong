package heartbeat

import (
	"object_storage_dong/lib/rabbitmq"
	"os"
	"time"
)

// 每隔5s向apiServers exchange发送一条消息
// 将本服务节点的监听对地址发送出去
func StartHeartbeat() {
	// 调用rabbitmq.New创建一个rabbitmq.RabbitMQ结构体
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	// 无限循环调用Publish方法向apiServers exchange发送本节点的监听地址
	for {
		q.Publish("apiServers", os.Getenv("LISTEN_ADDRESS"))
		time.Sleep(5 * time.Second)
	}
}
