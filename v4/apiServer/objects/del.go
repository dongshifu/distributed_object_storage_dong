package objects

import (
	"log"
	"net/http"
	"object_storage_dong/lib/es"
	"strings"
)

func del(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	// 以name为参数调用es.SearchLaestVersion,搜索该对象最新的版本
	version, e := es.SearchLatestVersion(name)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 插入新的元数据，接受元数据的name, version, size和hash
	// hash 为空字符串表示这个一个删除标记
	e = es.PutMetadata(name, version.Version+1, 0, "")
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
