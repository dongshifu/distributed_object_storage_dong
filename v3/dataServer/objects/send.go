package objects

import (
	"io"
	"os"
)

// 两个输入参数，用于写入对象数据的w和对象的文件名file
func sendFile(w io.Writer, file string) {
	// 调用os.Open打开对象文件
	f, _ := os.Open(file)
	defer f.Close()
	// 用io.Copy将文件内容写入w
	io.Copy(w, f)
}
