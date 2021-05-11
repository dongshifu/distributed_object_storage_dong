package main

import (
	"log"
	"net/http"
	"object_storage_dong/lib/es"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// 需要在每个数据节点上定期运行
	// 调用filepath.Global获取$STORAGE_ROOT/objects/目录下的所有文件
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")

	// for循环中遍历访问文件
	for i := range files {
		// 从文件中获得对象的散列值
		hash := strings.Split(filepath.Base(files[i]), ".")[0]
		// 调用es.HasHash检查元数据服务中是否存在该散列值
		hashInMetadata, e := es.HasHash(hash)
		if e != nil {
			log.Println(e)
			return
		}
		// 不存在，调用del删除散列值
		if !hashInMetadata {
			del(hash)
		}
	}
}

// 访问数据服务的DELETE对象接口进行散列值的删除
func del(hash string) {
	log.Println("delete", hash)
	url := "http://" + os.Getenv("LISTEN_ADDRESS") + "/objects/" + hash
	request, _ := http.NewRequest("DELETE", url, nil)
	client := http.Client{}
	client.Do(request)
}
