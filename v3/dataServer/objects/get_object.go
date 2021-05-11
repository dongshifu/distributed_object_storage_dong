package objects

import (
	"log"
	"net/url"
	"object_storage_dong/lib/utils"
	"object_storage_dong/v3/dataServer/locate"
	"os"
)

func getFile(hash string) string {
	// 先根据hash值找到对应的文件
	file := os.Getenv("STORAGE_ROOT") + "/objects/" + hash
	// 打开文件并计算文件的hash值和URL中的hash进行比较
	f, _ := os.Open(file)
	d := url.PathEscape(utils.CalculateHash(f))
	f.Close()
	// hash不一致，出现问题，从缓存和磁盘上删除对象
	// 返回空字符串
	if d != hash {
		log.Println("object hash mismatch, remove", file)
		locate.Del(hash)
		os.Remove(file)
		return ""
	}
	// 一致则返回对象的文件名
	return file
}
