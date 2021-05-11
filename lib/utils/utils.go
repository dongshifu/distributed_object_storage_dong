package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// Range的头部格式必须为bytes=<first>- 开头
func GetOffsetFromHeader(h http.Header) int64 {
	byteRange := h.Get("range")
	if len(byteRange) < 7 {
		return 0
	}
	if byteRange[:6] != "bytes=" {
		return 0
	}
	// 调用strings.Split将<first>部分切取出来
	bytePos := strings.Split(byteRange[6:], "-")
	// 调用strconv.ParseInt将字符串转换为int64返回
	offset, _ := strconv.ParseInt(bytePos[0], 0, 64)
	return offset
}

func GetHashFromHeader(h http.Header) string {
	// 获取"digest"头部
	digest := h.Get("digest")
	// 检查diest头部的形式是否满足"SHA-256=<hash>"
	if len(digest) < 9 {
		return ""
	}
	if digest[:8] != "SHA-256=" {
		return ""
	}
	return digest[8:]
}

func GetSizeFromHeader(h http.Header) int64 {
	// 得到"conten-length"头部，并调用strconv.PareseInt将字符串转换为int64输出
	size, _ := strconv.ParseInt(h.Get("content-length"), 0, 64)
	return size
}

func CalculateHash(r io.Reader) string {
	h := sha256.New()
	io.Copy(h, r)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
