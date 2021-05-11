package main

import (
	"log"
	"net/http"
	"object_storage_dong/v4/apiServer/heartbeat"
	"object_storage_dong/v4/apiServer/locate"
	"object_storage_dong/v4/apiServer/objects"
	"object_storage_dong/v4/apiServer/versions"
	"os"
)

func main() {
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", versions.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
