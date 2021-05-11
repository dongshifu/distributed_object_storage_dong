package objects

import (
	"fmt"
	"io"
	"object_storage_dong/lib/objectstream"
	"object_storage_dong/v3/apiServer/locate"
)

func getStream(object string) (io.Reader, error) {
	// 定位object对象
	server := locate.Locate(object)
	if server == "" {
		return nil, fmt.Errorf("object %s locate fail", object)
	}
	// 调用objectstream.NewGetStream并返回结果
	return objectstream.NewGetStream(server, object)
}
