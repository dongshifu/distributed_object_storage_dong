package temp

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type tempInfo struct {
	Uuid string
	Name string
	Size int64
}

func post(w http.ResponseWriter, r *http.Request) {
	// 生成一个随机的uuid
	output, _ := exec.Command("uuidgen").Output()
	uuid := strings.TrimSuffix(string(output), "\n")
	// 从请求的URL获取对象的名字，也即散列值
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	// 从头部获取对象的大小
	size, e := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 将uuid,name,size拼成一个tempInfo结构体
	t := tempInfo{uuid, name, size}
	// 调用tempInfo的writeToFile方法将结构体内容写入磁盘文件
	e = t.writeToFile()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 保存临时对象的内容
	os.Create(os.Getenv("STORAGE_ROOT") + "/temp/" + t.Uuid + ".dat")
	// 将uuid通过HTTP响应返回给接口服务
	w.Write([]byte(uuid))
}

func (t *tempInfo) writeToFile() error {
	f, e := os.Create(os.Getenv("STORAGE_ROOT") + "/temp/" + t.Uuid)
	if e != nil {
		return e
	}
	defer f.Close()
	// 将tempInfo的内容经过JSON编码后写入文件
	// 用于保存临时对象信息，与实际的对象内容不同
	b, _ := json.Marshal(t)
	f.Write(b)
	return nil
}
