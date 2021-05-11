package temp

import (
	"fmt"
	"log"
	"net/http"
	"object_storage_dong/lib/rs"
	"strings"
)

// 根据token恢复出stream后调用CurrentSize获取当前大小并方在Content-Length头部返回
func head(w http.ResponseWriter, r *http.Request) {
	token := strings.Split(r.URL.EscapedPath(), "/")[2]
	stream, e := rs.NewRSResumablePutStreamFromToken(token)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	current := stream.CurrentSize()
	if current == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("content-length", fmt.Sprintf("%d", current))
}
