package heartbeat

import (
	"math/rand"
)

// 从当前所有的数据服务节点中随机选择一个节点并返回
func ChooseRandomDataServer() string {
	ds := GetDataServers()
	n := len(ds)
	if n == 0 {
		return ""
	}
	return ds[rand.Intn(n)]
}
