package main

import (
	"log"
	"net/http"
	"object_storage_dong/v1/dataServer/objects"
	"object_storage_dong/v1/dataServer/heartbeat"
	"object_storage_dong/v1/dataServer/locate"
	"os"
)

func main() {
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
