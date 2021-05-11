package versions

import (
	"encoding/json"
	"log"
	"net/http"
	"object_storage_dong/lib/es"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// 检查HTTP方法是否为GET
	m := r.Method
	// 如果不为GET，返回405 Method Not Allowed
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// 方法为GET,获取URL中的<object_name>部分
	from := 0
	size := 1000
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	// 无限循环调用es包的SearchAllVersions函数
	for {
		// 返回一个元数据的数组
		metas, e := es.SearchAllVersions(name, from, size)
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// 遍历数组，将元数据一一写入HTTP响应的正文
		for i := range metas {
			b, _ := json.Marshal(metas[i])
			w.Write(b)
			w.Write([]byte("\n"))
		}
		// 如果返回的数组长度不等于size，说明元数据服务种没有更多的数据，直接返回
		if len(metas) != size {
			return
		}
		// 否则把from的值增加1000进行下一轮迭代
		from += size
	}
}
