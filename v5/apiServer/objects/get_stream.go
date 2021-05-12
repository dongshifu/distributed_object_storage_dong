package objects

import (
	"fmt"
	"object_storage_dong/lib/rs"
	"object_storage_dong/v5/apiServer/heartbeat"
	"object_storage_dong/v5/apiServer/locate"
)

// 提供给外部包使用。增加size参数是由于RS码的实现要求每一个数据片的长度完全一样
// 在编码如果对象长度不能被4整除，函数会队最后一个数片进行填充。
// 解码时必须提供对象的准确长度，防止填充数据被当成元数对象数据返回。
func GetStream(hash string, size int64) (*rs.RSGetStream, error) {
	// 根据对象散列值hash定位对象
	locateInfo := locate.Locate(hash)
	// 反馈的定位结果数组长度小于4，返回错误
	if len(locateInfo) < rs.DATA_SHARDS {
		return nil, fmt.Errorf("object %s locate fail, result %v", hash, locateInfo)
	}
	dataServers := make([]string, 0)
	// 长度不为6，说明对象有部分分片丢失
	if len(locateInfo) != rs.ALL_SHARDS {
		// 调用heartbeat.ChooseRandomDataServers随机选择用于接收恢复分片的数据服务节点
		dataServers = heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS-len(locateInfo), locateInfo)
	}
	// 调用rs.NewRSGetStream函数创建rs.RSGetStream
	return rs.NewRSGetStream(locateInfo, dataServers, hash, size)
}
