package objects

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.EscapedPath())
	//r.URL.EsccapedPath得到request的路径，此处为/objects/xxx
	f, e := os.Create(os.Getenv("STORAGE_ROOT") + "/objects/" +
		strings.Split(r.URL.EscapedPath(), "/")[2]) //得到文件名
	if e != nil {
		//创建文件失败
		log.Println(e)
		//写入HTTP响应的代码
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	io.Copy(f, r.Body) //将r.Body写入文件
}
