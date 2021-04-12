package objects

import (
	"log"
	"net/http"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	// 从URL中获取objects_name
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	// 将r.Body和objects作为参数调用storeObject
	// 第一个返回值为int类型的变量，用于表示HTTP错误码
	// 第二个返回值为error,如果error不为nil，出错并打印
	c, e := storeObject(r.Body, object)
	if e != nil {
		log.Println(e)
	}
	w.WriteHeader(c)
}
