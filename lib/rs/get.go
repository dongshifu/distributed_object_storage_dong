package rs

import (
	"fmt"
	"io"
	"object_storage_dong/lib/objectstream"
)

type RSGetStream struct {
	*decoder
}

func NewRSGetStream(locateInfo map[int]string, dataServers []string, hash string, size int64) (*RSGetStream, error) {
	// 检查是否满足4+2 的RS码需求
	// 不满足，返回错误
	if len(locateInfo)+len(dataServers) != ALL_SHARDS {
		return nil, fmt.Errorf("dataServers number mismatch")
	}

	// 创建长度为6的io.Reader数组readers，用于读取6个分片的数据
	readers := make([]io.Reader, ALL_SHARDS)
	// 遍历6个分片的id
	for i := 0; i < ALL_SHARDS; i++ {
		// 在locateinfo中查找分片所在的数据服务节点地址
		server := locateInfo[i]
		// 如果某个分片id相对的数据服务节点地址为空，说明该分片丢失
		// 取一个随机数据服务节点补上
		if server == "" {
			locateInfo[i] = dataServers[0]
			dataServers = dataServers[1:]
			continue
		}
		// 数据服务节点存在，调用objectstream.NewGetStream打开一个对象读取流用于读取该分片数据
		reader, e := objectstream.NewGetStream(server, fmt.Sprintf("%s.%d", hash, i))
		// 打开的流保存在readers数组相应的元素中
		if e == nil {
			readers[i] = reader
		}
	}

	writers := make([]io.Writer, ALL_SHARDS)
	perShard := (size + DATA_SHARDS - 1) / DATA_SHARDS
	var e error
	// 遍历readers
	for i := range readers {
		// 第一个次遍历出现nil的情况
		// 1.该分片数据服务节点地址为空
		// 2.数据服务节点存在但打开流失败
		if readers[i] == nil {
			// 某个元素为nil,调用NewTempPutStream创建相应的临时对象写入流用于恢复分片
			// 打开的流保存到writers数组相应元素中
			writers[i], e = objectstream.NewTempPutStream(locateInfo[i], fmt.Sprintf("%s.%d", hash, i), perShard)
			if e != nil {
				return nil, e
			}
		}
	}
	// 处理完成，readers和writers数组形成互补关系
	// 对于某个分片id，要么在readers中存在相应的读取流，要么在writers中存在相应的写入流
	// 将两个数组以及对象的大小size作为参数调用NewDecoder生成decoder结构体的指针dec
	dec := NewDecoder(readers, writers, size)
	// 将dec作为RSGetStream的内嵌结构体返回
	return &RSGetStream{dec}, nil
}

func (s *RSGetStream) Close() {
	for i := range s.writers {
		if s.writers[i] != nil {
			s.writers[i].(*objectstream.TempPutStream).Commit(true)
		}
	}
}

// offset表示需要跳过的字节数，whence表示起跳点
func (s *RSGetStream) Seek(offset int64, whence int) (int64, error) {
	// 方法只支持从当前位置io.SeekCurrent起跳
	if whence != io.SeekCurrent {
		panic("only support SeekCurrent")
	}
	// 跳过的字节数不能为负
	if offset < 0 {
		panic("only support forward seek")
	}
	// for循环中每次读取32000字节并丢弃，直到读到offset位置为止
	for offset != 0 {
		length := int64(BLOCK_SIZE)
		if offset < length {
			length = offset
		}
		buf := make([]byte, length)
		io.ReadFull(s, buf)
		offset -= length
	}
	return offset, nil
}
