package objects

import (
	"log"
	"net/http"
	"net/url"
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

	// 以散列值作为参数调用stroreObject
	c, e := storeObject(r.Body, url.PathEscape(hash))
	if e != nil {
		log.Println(e)
		w.WriteHeader(c)
		return
	}
	if c != http.StatusOK {
		w.WriteHeader(c)
		return
	}

	// 从URL中获取对象的名字和对象的大小
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	size := utils.GetSizeFromHeader(r.Header)
	// 以对象的名字、散列值和大小为参数调用es.AddVersions给对象添加新版本
	e = es.AddVersion(name, hash, size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
