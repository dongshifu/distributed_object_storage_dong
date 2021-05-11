package locate

import (
	"object_storage_dong/lib/rabbitmq"
	"object_storage_dong/lib/types"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var objects = make(map[string]int)
var mutex sync.Mutex

// 告知某个对象是否存在，同时告知本节点保存的是该对象哪个分片
func Locate(hash string) int {
	mutex.Lock()
	id, ok := objects[hash]
	mutex.Unlock()
	if !ok {
		return -1
	}
	return id
}

// 将对象及其分片id加入缓存
func Add(hash string, id int) {
	mutex.Lock()
	objects[hash] = id
	mutex.Unlock()
}

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
		// 读取来自接口服务需要定位的对象散列值hash
		hash, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		// 调用Locate获得分片id
		id := Locate(hash)
		if id != -1 {
			// id不为-1,将自身的节点监听地址和id打包成一个types.Locate Message结构体作为反馈消息发送
			q.Send(msg.ReplyTo, types.LocateMessage{Addr: os.Getenv("LISTEN_ADDRESS"), Id: id})
		}
	}
}

func CollectObjects() {
	// 获取存储目录下的所有文件
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
	for i := range files {
		// 以'.'分割其基本文件名，获得对象的散列值hash以及分片id
		file := strings.Split(filepath.Base(files[i]), ".")
		if len(file) != 3 {
			panic(files[i])
		}
		hash := file[0]
		id, e := strconv.Atoi(file[1])
		if e != nil {
			panic(e)
		}
		// 将对象的散列值hash以及分片id加入定位缓存
		objects[hash] = id
	}
}
