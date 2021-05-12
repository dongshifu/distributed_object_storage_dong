package main

import (
	"log"
	"object_storage_dong/lib/es"
	"object_storage_dong/lib/utils"
	"object_storage_dong/v7/apiServer/objects"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// 在数据服务节点上定期运行
	// 调用filepath.Glob获取$STOTAGE_ROOT/objects/目录下所有文件
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")

	// 在for循环中遍历访问文件
	for i := range files {
		// 从文件名中获得对象的散列值
		hash := strings.Split(filepath.Base(files[i]), ".")[0]
		// 调用verify检查数据
		verify(hash)
	}
}

func verify(hash string) {
	log.Println("verify", hash)
	// 调用es.SearchHashSize从元数据服务中获取该散列值对应的对象大小
	size, e := es.SearchHashSize(hash)
	if e != nil {
		log.Println(e)
		return
	}
	// 以对象的散列值和大小为参数调用objects.GetStream创建一个对象数据流
	// 底层实现会自动完成数据的修复。
	stream, e := objects.GetStream(hash, size)
	if e != nil {
		log.Println(e)
		return
	}
	// 调用utils.CalculateHash计算对象的散列值
	d := utils.CalculateHash(stream)
	// 检查hash是否一致，不一致则以log的形式报告错误(数据损坏，已经不可修复)
	if d != hash {
		log.Printf("object hash mismatch, calculated=%s, requested=%s", d, hash)
	}
	// 关闭数据对象流
	stream.Close()
}
