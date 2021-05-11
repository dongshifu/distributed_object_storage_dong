package rs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"object_storage_dong/lib/objectstream"
	"object_storage_dong/lib/utils"
)

// 保存对象的名字、大小、散列值以及6个分片所在数据服务节点的地址和uuid
type resumableToken struct {
	Name    string
	Size    int64
	Hash    string
	Servers []string
	Uuids   []string
}

type RSResumablePutStream struct {
	*RSPutStream
	*resumableToken
}

// 传入保存数据服务节点地址的dataServers数组，对象的名字name，散列值以及大小
func NewRSResumablePutStream(dataServers []string, name, hash string, size int64) (*RSResumablePutStream, error) {
	// 调用NewRSPutStream创建一个类型为RSPutStream的变量putStream
	putStream, e := NewRSPutStream(dataServers, hash, size)
	if e != nil {
		return nil, e
	}
	uuids := make([]string, ALL_SHARDS)
	// 从putStream的成员writers数组中获取6个分片的uuid，保存到uuid数组
	for i := range uuids {
		uuids[i] = putStream.writers[i].(*objectstream.TempPutStream).Uuid
	}
	// 创建resumableToken结构体token
	token := &resumableToken{name, size, hash, dataServers, uuids}
	return &RSResumablePutStream{putStream, token}, nil
}

func NewRSResumablePutStreamFromToken(token string) (*RSResumablePutStream, error) {
	// 对token进行Base64解码
	b, e := base64.StdEncoding.DecodeString(token)
	if e != nil {
		return nil, e
	}

	var t resumableToken
	// 将JSON数据编出形成resumableToken结构体t
	e = json.Unmarshal(b, &t)
	if e != nil {
		return nil, e
	}

	writers := make([]io.Writer, ALL_SHARDS)
	// t的Servers和uuid数组中保存了当初创建的6个分片临时对象所在的数据服务节点地址和uuid
	for i := range writers {
		// 保存到writers数组中
		writers[i] = &objectstream.TempPutStream{t.Servers[i], t.Uuids[i]}
	}
	// 以writers数组为参数创建encoder结构体enc
	enc := NewEncoder(writers)
	// 以enc为内嵌结构体创建RSPutStream
	// 最终以RSPutStream和t为内嵌结构体创建RSResumablePutStream返回
	return &RSResumablePutStream{&RSPutStream{enc}, &t}, nil
}

// 将自身数据以JSON格式编入，然后返回结果Base64编码后的字符串
func (s *RSResumablePutStream) ToToken() string {
	b, _ := json.Marshal(s)
	return base64.StdEncoding.EncodeToString(b)
}

// 以HEAD方法获取第一个分片临时对象的大小并乘4作为size返回
func (s *RSResumablePutStream) CurrentSize() int64 {
	r, e := http.Head(fmt.Sprintf("http://%s/temp/%s", s.Servers[0], s.Uuids[0]))
	if e != nil {
		log.Println(e)
		return -1
	}
	if r.StatusCode != http.StatusOK {
		log.Println(r.StatusCode)
		return -1
	}
	size := utils.GetSizeFromHeader(r.Header) * DATA_SHARDS
	// 如果size超出对象的大小，返回对象大小
	if size > s.Size {
		size = s.Size
	}
	return size
}
