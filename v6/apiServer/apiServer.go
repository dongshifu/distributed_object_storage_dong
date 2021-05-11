package main

import (
	"log"
	"net/http"
	"object_storage_dong/v6/apiServer/heartbeat"
	"object_storage_dong/v6/apiServer/locate"
	"object_storage_dong/v6/apiServer/objects"
	"object_storage_dong/v6/apiServer/temp"
	"object_storage_dong/v6/apiServer/versions"
	"os"
)

func main() {
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	// 增加temp.Handler函数用于处理对/temp的请求
	http.HandleFunc("/temp/", temp.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", versions.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
