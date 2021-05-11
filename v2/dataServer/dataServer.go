package main

import (
	"log"
	"net/http"
	"object_storage_dong/v2/dataServer/heartbeat"
	"object_storage_dong/v2/dataServer/locate"
	"object_storage_dong/v2/dataServer/objects"
	"os"
)

func main() {
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
