package heartbeat

import (
	"math/rand"
)

// 用于获取上传复原分片的随机数据服务节点。目前已有的分片所在的数据服务节点需要被排除
func ChooseRandomDataServers(n int, exclude map[int]string) (ds []string) {
	// 输入参数n表示需要多少个随机数服务节点
	// 输入参数exclude表明要求返回的数据服务节点不能包含哪些节点。
	// 当定位完成后，实际收到的反馈消息可能不足6个，此时需要进行数据恢复，即根据目前已有的分片将丢失的分片复原出来并再次上传到数据服务
	candidates := make([]string, 0)
	reverseExcludeMap := make(map[string]int)
	for id, addr := range exclude {
		reverseExcludeMap[addr] = id
	}
	servers := GetDataServers()
	for i := range servers {
		s := servers[i]
		_, excluded := reverseExcludeMap[s]
		if !excluded {
			// 不需要被排除的加入到condaidate数组
			candidates = append(candidates, s)
		}
	}
	length := len(candidates)
	if length < n {
		// 无法满足要求的n个数据服务节点，返回一个空数组
		return
	}
	// 将0-length-1的所有整数乱序排列返回一个数组
	p := rand.Perm(length)
	// 取前n个作为candicate数组的下标取数据节点地址返回
	for i := 0; i < n; i++ {
		ds = append(ds, candidates[p[i]])
	}
	return
}
