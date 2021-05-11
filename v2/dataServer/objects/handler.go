package objects

import "net/http"

// 检查HTTP请求方法:PUT则调用put函数，GET则调用共get函数。其余则返回405 Method Not Allowed错误代码
func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method //Method记录该HTTP请求的方法
	if m == http.MethodPut {
		put(w, r)
		return
	}
	if m == http.MethodGet {
		get(w, r)
		return
	}
	//写HTTP响应的代码
	w.WriteHeader(http.StatusMethodNotAllowed)
}
