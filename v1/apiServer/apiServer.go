package main

import (
	"log"
	"net/http"
	"object_storage_dong/v1/apiServer/heartbeat"
	"object_storage_dong/v1/apiServer/locate"
	"object_storage_dong/v1/apiServer/objects"
	"os"
)

func main() {
	// 提供locate功能
	go heartbeat.ListenHeartbeat()
	// 处理URL以/objects/开头的对象请求
	http.HandleFunc("/objects/", objects.Handler)
	// 处理URL以/locate/开头的定位请求
	http.HandleFunc("/locate/", locate.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
