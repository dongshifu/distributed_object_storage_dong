package objects

import (
	"log"
	"net/http"
	"net/url"
	"object_storage_dong/lib/es"
	"object_storage_dong/lib/rs"
	"object_storage_dong/lib/utils"
	"object_storage_dong/v5/apiServer/heartbeat"
	"object_storage_dong/v5/apiServer/locate"
	"strconv"
	"strings"
)

func post(w http.ResponseWriter, r *http.Request) {
	// 从请求的URL中获得对象的名字
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	// 从请求的相应头部获得对象的大小
	size, e := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// 从请求的相应头部获得对象的散列值
	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Println("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 如果散列值已经存在，可以直接往元数据服务添加新版本并返回200OK
	if locate.Exist(url.PathEscape(hash)) {
		e = es.AddVersion(name, hash, size)
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		return
	}
	// 散列值不存在，随机选出6个数据节点
	ds := heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS, nil)
	if len(ds) != rs.ALL_SHARDS {
		log.Println("cannot find enough dataServer")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	// 调用rs.NewRSResumablePutStream生成数据流stream
	stream, e := rs.NewRSResumablePutStream(ds, name, url.PathEscape(hash), size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 调用ToToken方法生成一个字符串token,放入Location响应头部
	w.Header().Set("location", "/temp/"+url.PathEscape(stream.ToToken()))
	// 返回HTTP代码201 Created
	w.WriteHeader(http.StatusCreated)
}
