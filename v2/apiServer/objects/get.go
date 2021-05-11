package objects

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"object_storage_dong/lib/es"
	"strconv"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	// 获取URL并从URL的查询参数中获取"version"参数的值
	// Query方法返回一个保存URL所有查询参数的map，该map的键是查询参数的名字，值是一个字符串数组
	// HTTP的URL查询参数允许存在多个值，以"version"为key可以得到URL中查询参数的所有值
	versionId := r.URL.Query()["version"]
	version := 0
	var e error
	if len(versionId) != 0 {
		// 项目中不考虑多个"version"查询参数的情况
		// 始终以versionId数组的第一个元素作为客户端提供的版本号
		// 将字符串转换为整型
		version, e = strconv.Atoi(versionId[0])
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	// 调用es的GetMetadata函数得到对象的元数据meta
	meta, e := es.GetMetadata(name, version)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// meta.Hash为对象的散列值，如果为空表示该对象版本是一个删除标记
	// 返回404 Not Found
	if meta.Hash == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// 以散列值为对象名从数据服务层获取对象并输出
	object := url.PathEscape(meta.Hash)
	stream, e := getStream(object)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	io.Copy(w, stream)
}
