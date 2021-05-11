package locate

import (
	"encoding/json"
	"net/http"
	"strings"
)

// 处理HTTP请求
func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// 将文件名作为Locate函数的参数进行定位
	info := Locate(strings.Split(r.URL.EscapedPath(), "/")[2])
	// 为空，说明定位失败
	if len(info) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// 不为空，则拥有该对象的一个数据服务节点的地址，将地址作为HTTP响应的正文输出
	b, _ := json.Marshal(info)
	w.Write(b)
}
