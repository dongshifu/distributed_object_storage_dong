package objects

import (
	"fmt"
	"object_storage_dong/lib/rs"
	"object_storage_dong/v5/apiServer/heartbeat"
)

func putStream(hash string, size int64) (*rs.RSPutStream, error) {
	servers := heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS, nil)
	if len(servers) != rs.ALL_SHARDS {
		return nil, fmt.Errorf("cannot find enough dataServer")
	}

	// 返回一个指向rs.RSPutStream结构体的指针
	return rs.NewRSPutStream(servers, hash, size)
}
