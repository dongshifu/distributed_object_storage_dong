package objects

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	f, e := os.Open(os.Getenv("STORAGE_ROOT") + "/objects/" +
		strings.Split(r.URL.EscapedPath(), "/")[2])
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	io.Copy(w, f)
	//f本身的类型是*os.File,同时实现了io.Writer和io.Reader两个接口，即实现了Write和Read方法
	//http.ResponseWriter也是接口，该接口实现了Write方法，也是一个io.Write接口
}
