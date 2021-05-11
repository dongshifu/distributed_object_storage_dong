package objects

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"object_storage_dong/lib/es"
	"object_storage_dong/lib/utils"
	"strconv"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	versionId := r.URL.Query()["version"]
	version := 0
	var e error
	if len(versionId) != 0 {
		version, e = strconv.Atoi(versionId[0])
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	meta, e := es.GetMetadata(name, version)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if meta.Hash == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	hash := url.PathEscape(meta.Hash)
	stream, e := GetStream(hash, meta.Size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// 调用utils.GetOffsetFromHeader函数获取HTTP的Range头部
	offset := utils.GetOffsetFromHeader(r.Header)
	if offset != 0 {
		// offset不为0，调用stream的Seek方法跳到offset位置
		stream.Seek(offset, io.SeekCurrent)
		// 设置Content-Range响应头部以及HTTP代码206 Partial Content
		w.Header().Set("content-range", fmt.Sprintf("bytes %d-%d/%d", offset, meta.Size-1, meta.Size))
		w.WriteHeader(http.StatusPartialContent)
	}
	// 通过io.Copy输出数据
	io.Copy(w, stream)
	stream.Close()
}
