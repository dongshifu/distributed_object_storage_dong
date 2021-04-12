package objectstream

import (
	"fmt"
	"io"
	"net/http"
)

type PutStream struct {
	// writer用于实现Write方法
	writer *io.PipeWriter
	// c用于把在一个gouroutine传输数据过程中发生的错误传回主线程
	c      chan error
}

// 用于生成一个PutStream结构体
func NewPutStream(server, object string) *PutStream {
	// 用io.Pipe创建一对reader和writer,类型为*io.PipeReader和*io.PipeWriter
	// 管道互联，写入writer的内容可以从reader中读出
	// 希望以写入数据流的方法操作HTTP的PUT请求
	reader, writer := io.Pipe()
	c := make(chan error)
	go func() {
		// 生成put请求，需要提供一个io.Reader作为http.NewRequest的参数
		request, _ := http.NewRequest("PUT", "http://"+server+"/objects/"+object, reader)
		// http.Client负责从request中读取需要PUT的内容
		client := http.Client{}
		// 由于管道的读写阻塞特性，在goroutine中调用client.Do方法
		r, e := client.Do(request)
		if e == nil && r.StatusCode != http.StatusOK {
			e = fmt.Errorf("dataServer return http code %d", r.StatusCode)
		}
		c <- e
	}()
	return &PutStream{writer, c}
}

// 用于写入writer,实现该方法PutStream才被认为实现了io.Write接口
func (w *PutStream) Write(p []byte) (n int, err error) {
	return w.writer.Write(p)
}

// 关闭writer,为了让管道另一端的reader读到io.EOF，否则在gouroutine中运行的client.Do将始终阻塞无法返回
func (w *PutStream) Close() error {
	w.writer.Close()
	// 从c中读取发送自goroutine得错误并返回
	return <-w.c
}
