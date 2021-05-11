package objects

import (
	"compress/gzip"
	"io"
	"log"
	"os"
)

func sendFile(w io.Writer, file string) {
	f, e := os.Open(file)
	if e != nil {
		log.Println(e)
		return
	}
	defer f.Close()
	// 在对象文件上用gzip.NewReader创建一个指向gzip.Reader结构体的指针gzipStream
	gzipStream, e := gzip.NewReader(f)
	if e != nil {
		log.Println(e)
		return
	}
	// 读出gzipStream中的数据
	io.Copy(w, gzipStream)
	gzipStream.Close()
}
