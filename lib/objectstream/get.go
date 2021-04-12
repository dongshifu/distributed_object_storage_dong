package objectstream

import (
	"fmt"
	"io"
	"net/http"
)

type GetStream struct {
	reader io.Reader
}

func newGetStream(url string) (*GetStream, error) {
	// 输入的url表示用于获取数据流的HTTP服务地址
	// 调用http.Get发起一个GET请求，获取该地址的HTTP响应
	r, e := http.Get(url) //返回的r类型为*http.Response,其body是用于读取HTTP响应正文的io.Reader
	if e != nil {
		return nil, e
	}
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dataServer return http code %d", r.StatusCode)
	}
	return &GetStream{r.Body}, nil
}

// 封装newGetStream函数
func NewGetStream(server, object string) (*GetStream, error) {
	if server == "" || object == "" {
		return nil, fmt.Errorf("invalid server %s object %s", server, object)
	}
	// 内部拼凑一个url传给newGetStream，对外隐藏url的细节
	return newGetStream("http://" + server + "/objects/" + object)
}

// 用于读取reader成员，实现该方法，则GetStream结构体实现io.Reader接口
func (r *GetStream) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}
