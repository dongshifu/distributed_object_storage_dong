package rs

import (
	"io"

	"github.com/klauspost/reedsolomon"
)

type encoder struct {
	writers []io.Writer
	enc     reedsolomon.Encoder // reedsolomon.Encoder接口
	cache   []byte              //用于做输入数据缓存的字节数组
}

// 调用reeddolomon.New生成具有4个数据片加两个校验片的RS码编码器enc
func NewEncoder(writers []io.Writer) *encoder {
	enc, _ := reedsolomon.New(DATA_SHARDS, PARITY_SHARDS)
	// 将输入参数writers和enc作为生成的encoder结构体的成员返回
	return &encoder{writers, enc, nil}
}

// RSPutStream本身没有实现Write方法，实现的时候函数会直接调用其内嵌结构体encoder的Write方法
func (e *encoder) Write(p []byte) (n int, err error) {
	length := len(p)
	current := 0
	// 将p中待写入的数据以块的形式放入缓存
	for length != 0 {
		next := BLOCK_SIZE - len(e.cache)
		if next > length {
			next = length
		}
		e.cache = append(e.cache, p[current:current+next]...)
		// 如果缓存已满，调用Flush方法将缓存实际写入writers
		// 缓存的上限是每个数据片8000字节，4个数据片共32000字节。
		// 如果缓存里剩余的数据不满32000字节就暂不刷新，等待Write方法下一次调用
		if len(e.cache) == BLOCK_SIZE {
			e.Flush()
		}
		current += next
		length -= next
	}
	return len(p), nil
}

func (e *encoder) Flush() {
	if len(e.cache) == 0 {
		return
	}
	// 调用Split方法将缓存的数据切成4个数据片
	shards, _ := e.enc.Split(e.cache)
	// 调用Encode方法生成两个校验片
	e.enc.Encode(shards)
	// 循环将6个片的数据依次写入wtiters并清空缓存
	for i := range shards {
		e.writers[i].Write(shards[i])
	}
	e.cache = []byte{}
}
