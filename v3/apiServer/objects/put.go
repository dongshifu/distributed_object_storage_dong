package objects

import (
	"log"
	"net/http"
	"object_storage_dong/lib/es"
	"object_storage_dong/lib/utils"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	// 先从HTTP请求头部获取对象的散列值
	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Println("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 从URL中获取对象的大小
	size := utils.GetSizeFromHeader(r.Header)
	// 以散列值和size作为参数调用stroreObject
	// 新实现的storeObject需要在一开始就确定临时对象大小
	c, e := storeObject(r.Body, hash, size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(c)
		return
	}
	if c != http.StatusOK {
		w.WriteHeader(c)
		return
	}

	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	e = es.AddVersion(name, hash, size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
