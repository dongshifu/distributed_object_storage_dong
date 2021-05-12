package objects

import (
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net/url"
	"object_storage_dong/v7/dataServer/locate"
	"os"
	"path/filepath"
	"strings"
)

func getFile(name string) string {
	// 在$STORAGE_ROOT/objects/目录下查找所有以<hash>.X开头的文件
	files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/" + name + ".*")
	// 找不到，返回空字符串
	if len(files) != 1 {
		return ""
	}
	// 找到之后计算散列值
	file := files[0]
	h := sha256.New()
	sendFile(h, file)
	d := url.PathEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	hash := strings.Split(file, ".")[2]
	// 如果与<hash of shard X>的值不匹配则删除该对象并返回空字符串
	if d != hash {
		log.Println("object hash mismatch, remove", file)
		locate.Del(hash)
		os.Remove(file)
		return ""
	}
	// 否则返回对象的文件名
	return file
}
