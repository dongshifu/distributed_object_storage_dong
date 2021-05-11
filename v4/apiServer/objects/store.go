package objects

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"object_storage_dong/lib/utils"
	"object_storage_dong/v4/apiServer/locate"
)

func storeObject(r io.Reader, hash string, size int64) (int, error) {
	// 首先调用locate.Exist方法定位对象的散列值
	// 如果已经存在，跳过后续上传操作直接返回200 OK
	if locate.Exist(url.PathEscape(hash)) {
		return http.StatusOK, nil
	}

	// 不存在，调用putStream生成对象的写入流stream用于写入
	stream, e := putStream(url.PathEscape(hash), size)
	if e != nil {
		return http.StatusInternalServerError, e
	}

	// 两个输入参数，分别是作为io.Reader的r和io.Writer的stream
	// 返回的reader也是一个io.Reader
	reader := io.TeeReader(r, stream)
	// reader被读取的时候，实际的内容读取自r,同时也会写入stream
	// 用utils.CalculateHash从reader中读取数据的同时也写入了stream
	d := utils.CalculateHash(reader)
	// 计算出来的散列值和hash做比较
	// 不一致则调用stream.Commit(false)删除临时对象，并返回400 Bad Request
	if d != hash {
		stream.Commit(false)
		return http.StatusBadRequest, fmt.Errorf("object hash mismatch, calculated=%s, requested=%s", d, hash)
	}
	// 一致则调用stream.Commit(true)将临时对象转正并返回200 OK
	stream.Commit(true)
	return http.StatusOK, nil
}
