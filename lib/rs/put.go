package rs

import (
	"fmt"
	"io"
	"object_storage_dong/lib/objectstream"
)

// 内嵌一个encoder结构体，RPutStream的使用者可以像访问RPutStream的方法或成员一样访问*encoder的方法或成员
type RSPutStream struct {
	*encoder
}

func NewRSPutStream(dataServers []string, hash string, size int64) (*RSPutStream, error) {
	// dataServers长度不为6,返回错误
	if len(dataServers) != ALL_SHARDS {
		return nil, fmt.Errorf("dataServers number mismatch")
	}
	// 根据size计算出每个分片的大小perShard,也就是size/4再向上取整
	perShard := (size + DATA_SHARDS - 1) / DATA_SHARDS
	// 创建长度为6的io.Writers数组
	writers := make([]io.Writer, ALL_SHARDS)
	var e error
	for i := range writers {
		// writers数组中的每个元素都是一个objectstream.TempPutSterm,用于上传一个分片对象
		writers[i], e = objectstream.NewTempPutStream(dataServers[i],
			fmt.Sprintf("%s.%d", hash, i), perShard)
		if e != nil {
			return nil, e
		}
	}
	// 调用NewEncoder函数创建一个encoder结构体指针enc
	enc := NewEncoder(writers)
	// 将enc作为RSPutSterm的内嵌结构体返回
	return &RSPutStream{enc}, nil
}

func (s *RSPutStream) Commit(success bool) {
	// 调用RSPutStream的内嵌结构体encoder的Flush方法将缓存中最后的数据写入
	s.Flush()
	// 对调用encoder的成员数组writers中的元素调用Commit方法将6个临时对象依次转正或删除
	for i := range s.writers {
		s.writers[i].(*objectstream.TempPutStream).Commit(success)
	}
}
