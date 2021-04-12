package objects

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	// 调用getStream生成一个类型为io.Reader的stream
	stream, e := getStream(object)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// 调用io.Copy将stream的内容写入HTTP响应的正文
	io.Copy(w, stream)
}
