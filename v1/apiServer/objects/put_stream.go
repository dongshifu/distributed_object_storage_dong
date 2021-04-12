package objects

import (
	"fmt"
	"object_storage_dong/lib/objectstream"
	"object_storage_dong/v1/apiServer/heartbeat"
)

func putStream(object string) (*objectstream.PutStream, error) {
	// 先获得一个随机数服务节点的地址
	server := heartbeat.ChooseRandomDataServer()
	fmt.Println("data server =", server)
	// 没有可用的数据服务节点，返回objectstream.PutStream的空指针
	if server == "" {
		return nil, fmt.Errorf("cannot find any dataServer")
	}

	return objectstream.NewPutStream(server, object), nil
}
