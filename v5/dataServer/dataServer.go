package main

import (
	"log"
	"net/http"
	"object_storage_dong/v5/dataServer/heartbeat"
	"object_storage_dong/v5/dataServer/locate"
	"object_storage_dong/v5/dataServer/objects"
	"object_storage_dong/v5/dataServer/temp"
	"os"
)

func main() {
	// 之前的版本定位对象通过调用os.Stat来检查对象文件是否存在
	// 每次定位请求都会导致一次磁盘访问，会对系统带来很大负担
	// 为减少对磁盘的访问次数，数据服务定位功能仅在程序启动时候扫描一遍本地磁盘
	// 将磁盘中所有的对象散列值读入内存，之后的定位不需要再次访问磁盘，只需搜索内存即可
	locate.CollectObjects()
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler)
	// 引入temp.Handler处理函数注册
	http.HandleFunc("/temp/", temp.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
