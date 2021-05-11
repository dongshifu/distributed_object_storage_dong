package objects

import (
	"io"
	"net/http"
)

func storeObject(r io.Reader, object string) (int, error) {
	stream, e := putStream(object)
	// 没有找到可用的数据服务节点
	if e != nil {
		return http.StatusServiceUnavailable, e
	}

	// 找到可用的数服务节点并得到一个objectstream.PutStream的指针
	// objectstream.PutStream实现了Write方法，是一个io.Write接口
	// 用io.Copy将HTTP请求的正文写入stream
	io.Copy(stream, r)
	e = stream.Close()
	//写入出错
	if e != nil {
		return http.StatusInternalServerError, e
	}
	//写入成功
	return http.StatusOK, nil
}
