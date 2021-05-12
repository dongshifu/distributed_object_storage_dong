package locate

import (
	"encoding/json"
	"object_storage_dong/lib/rabbitmq"
	"object_storage_dong/lib/rs"
	"object_storage_dong/lib/types"
	"os"
	"time"
)

func Locate(name string) (locateInfo map[int]string) {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	q.Publish("dataServers", name)
	c := q.Consume()
	go func() {
		// 设置1秒超时，无论当前收到多少条反馈消息都会立即返回
		time.Sleep(time.Second)
		q.Close()
	}()
	locateInfo = make(map[int]string)
	// 使用for循环获取6条消息，每条消息包含拥有某个分片的数据服务节点的地址和分片的id
	// rs.ALL_SHARDS为rs包中的常数6,代表一共有4+2个分片
	for i := 0; i < rs.ALL_SHARDS; i++ {
		msg := <-c
		if len(msg.Body) == 0 {
			return
		}
		var info types.LocateMessage
		json.Unmarshal(msg.Body, &info)
		// 获取到的分片信息放到输出参数locateInfo中
		locateInfo[info.Id] = info.Addr
	}
	return
}

// 判断收到的反馈消息数量是否大于等于4
func Exist(name string) bool {
	// 大于等于4则说明对象存在，否则说明对象不存在(或者存在也无法读取)
	return len(Locate(name)) >= rs.DATA_SHARDS
}
