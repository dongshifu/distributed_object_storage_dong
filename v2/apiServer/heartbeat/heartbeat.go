package heartbeat

import (
	"object_storage_dong/lib/rabbitmq"
	"os"
	"strconv"
	"sync"
	"time"
)

// map,整个包可见，用于缓存所有的数据服务节点
var dataServers = make(map[string]time.Time)
var mutex sync.Mutex

func ListenHeartbeat() {
	// 创建消息队列绑定apiServers exchange
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	q.Bind("apiServers")
	// 通过go channel监听每个来自数据服务节点的心跳信息
	c := q.Consume()
	go removeExpiredDataServer() //goroutine并行处理
	for msg := range c {
		//将数据服务节点的监听地址作为map的键，收到消息的时间作为值存入map中
		dataServer, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		mutex.Lock()
		dataServers[dataServer] = time.Now()
		mutex.Unlock()
	}
}

// 每隔5s扫描一遍map,清除其中超过10s没有收到心跳消息的数据服务节点
func removeExpiredDataServer() {
	for {
		time.Sleep(5 * time.Second)
		mutex.Lock()
		for s, t := range dataServers {
			if t.Add(10 * time.Second).Before(time.Now()) {
				delete(dataServers, s)
			}
		}
		mutex.Unlock()
	}
}

// 遍历map并返回当前所有的数据服务节点
// 为防止多个goroutine并发读写map造成错误，map读写全部需要mutex的保护
func GetDataServers() []string {
	mutex.Lock()
	defer mutex.Unlock()
	ds := make([]string, 0)
	for s, _ := range dataServers {
		ds = append(ds, s)
	}
	return ds
}
