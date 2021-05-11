package rs

import (
	"io"

	"github.com/klauspost/reedsolomon"
)

type decoder struct {
	readers   []io.Reader
	writers   []io.Writer
	enc       reedsolomon.Encoder // 接口，用于RS解码
	size      int64               // 对象的大小
	cache     []byte              // cache用于缓存读取的数据
	cacheSize int                 // cachSize
	total     int64               // 当前读了多少字节
}

func NewDecoder(readers []io.Reader, writers []io.Writer, size int64) *decoder {
	// 调用reedsolomon.New创建4+2 RS码的解码器enc
	enc, _ := reedsolomon.New(DATA_SHARDS, PARITY_SHARDS)
	// 设置decoder结构体中相应属性后返回
	return &decoder{readers, writers, enc, size, nil, 0, 0}
}

func (d *decoder) Read(p []byte) (n int, err error) {
	// 当cache中没有更多数据时会调用getData方法获取数据
	if d.cacheSize == 0 {
		// getData返回的e不为nil，说明没能获取更多数据
		e := d.getData()
		if e != nil {
			// 返回0和e通知调用方
			return 0, e
		}
	}
	// 函数参数p的数组长度
	length := len(p)
	if d.cacheSize < length {
		// length超出当前缓冲区的数据大小，令length等于缓存的数据大小
		length = d.cacheSize
	}
	// 将缓存中length长度的数据复制给输入参数p，然后调整缓存，仅保留剩下的部分
	d.cacheSize -= length
	copy(p, d.cache[:length])
	d.cache = d.cache[length:]
	// 返回length,通知调用方本次读取一共有多少数据被复制到p中
	return length, nil
}

func (d *decoder) getData() error {
	// 先判断当前解码的数据大小是否等于对象原始大小
	// 如果已经相等，说明所有数都已经被读取，返回io.EOF
	if d.total == d.size {
		return io.EOF
	}
	// 如果还有数需要读取，创建一个长度为6的数组
	// 每个元素都是一个字节数组，用于保存相应分片中读取的数据
	shards := make([][]byte, ALL_SHARDS)
	repairIds := make([]int, 0)
	// 遍历6个shards
	for i := range shards {
		// 若某个分片对应的reader为nil，说明该分片已经丢失
		if d.readers[i] == nil {
			// 在repairIds中添加该分片的id
			repairIds = append(repairIds, i)
		} else {
			// 对应的reader不为nil，那么对应的shards需要被初始化为一个长度为8000的字节数组
			shards[i] = make([]byte, BLOCK_PER_SHARD)
			// 调用io.ReadFull从reader中完整读取8000字节的数据保存在shards中。
			n, e := io.ReadFull(d.readers[i], shards[i])
			// 如果发生了非EOF失败，该shards被置为nil
			if e != nil && e != io.EOF && e != io.ErrUnexpectedEOF {
				shards[i] = nil
				// 读取数据长度n不到8000字节，将shards实际的长度缩减为n
			} else if n != BLOCK_PER_SHARD {
				shards[i] = shards[i][:n]
			}
		}
	}
	// 调用成员enc的Reconstruct方法尝试将被置为nil的shards恢复出来
	e := d.enc.Reconstruct(shards)
	// 若发生错误，说明对象遭到不可恢复的破坏，返回错误给上层
	if e != nil {
		return e
	}
	// 恢复成功，6个shards中都保存了对应分片的正确数据
	// 遍历repairIds，将需要恢复的分片数据写入到相应的writer
	for i := range repairIds {
		id := repairIds[i]
		d.writers[id].Write(shards[id])
	}
	// 遍历4个数据分片
	for i := 0; i < DATA_SHARDS; i++ {
		shardSize := int64(len(shards[i]))
		if d.total+shardSize > d.size {
			shardSize -= d.total + shardSize - d.size
		}
		// 将每个分片中的数据添加到缓存cache中
		d.cache = append(d.cache, shards[i][:shardSize]...)
		// 修改缓存当前的大小cacheSize
		d.cacheSize += int(shardSize)
		// 修改当前已经读取的全部数据的大小total
		d.total += shardSize
	}
	return nil
}
