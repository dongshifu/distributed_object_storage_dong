package objects

import (
	"fmt"
	"object_storage_dong/lib/objectstream"
	"object_storage_dong/v3/apiServer/heartbeat"
)

func putStream(hash string, size int64) (*objectstream.TempPutStream, error) {
	server := heartbeat.ChooseRandomDataServer()
	if server == "" {
		return nil, fmt.Errorf("cannot find any dataServer")
	}

	// 数据服务的temp接口代替了原先的对象PUT接口，调用objectstream.NewTempPutStream
	return objectstream.NewTempPutStream(server, hash, size)
}
