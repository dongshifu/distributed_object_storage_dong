package locate

import (
	"object_storage_dong/lib/rabbitmq"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

// 用于缓存所有对象
var objects = make(map[string]int)
// 保护对objects的读写操作
var mutex sync.Mutex

// 利用map操作判断某个散列值是否存在于objects中，存在返回true,否则返回false
func Locate(hash string) bool {
	mutex.Lock()
	_, ok := objects[hash]
	mutex.Unlock()
	return ok
}

// 将一个散列值加入缓存，输入参数hash作为存入map的键，值为1
func Add(hash string) {
	mutex.Lock()
	objects[hash] = 1
	mutex.Unlock()
}

// 将一个散列值移出缓存
func Del(hash string) {
	mutex.Lock()
	delete(objects, hash)
	mutex.Unlock()
}

func StartLocate() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	q.Bind("dataServers")
	c := q.Consume()
	for msg := range c {
		hash, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}

		// 直接将从RabbitMQ消息队列中收到的对象散列值作为Locate参数
		exist := Locate(hash)
		if exist {
			q.Send(msg.ReplyTo, os.Getenv("LISTEN_ADDRESS"))
		}
	}
}

func CollectObjects() {
	// 读取存储目录里的所有文件
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
	for i := range files {
		// 对读出的文件一一调用filepath.Base获取其基本文件名
		// 也就是对象的散列值，将散列值加入objects缓存
		hash := filepath.Base(files[i])
		objects[hash] = 1
	}
}
