package temp

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"object_storage_dong/lib/es"
	"object_storage_dong/lib/rs"
	"object_storage_dong/lib/utils"
	"object_storage_dong/v6/apiServer/locate"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	token := strings.Split(r.URL.EscapedPath(), "/")[2]
	stream, e := rs.NewRSResumablePutStreamFromToken(token)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// 调用CurrentSize获取token当前大小
	current := stream.CurrentSize()
	if current == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// 从Range头部获得offset
	offset := utils.GetOffsetFromHeader(r.Header)
	// 如果offset和当前的大小不一致，接口服务返回416Range Not Satisfiable
	if current != offset {
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}
	bytes := make([]byte, rs.BLOCK_SIZE)
	// 在for循环中以32000字节为长度读取HTTP请求的正文并写入stream
	for {
		n, e := io.ReadFull(r.Body, bytes)
		if e != nil && e != io.EOF && e != io.ErrUnexpectedEOF {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		current += int64(n)
		// 如果读到的总长度超出了对象的大小，说明客户端上传的数据有误
		// 接口服务删除临时对象并返回403 Forbidden
		if current > stream.Size {
			stream.Commit(false)
			log.Println("resumable put exceed size")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if n != rs.BLOCK_SIZE && current != stream.Size {
			return
		}
		stream.Write(bytes[:n])
		// 读到的总长度等于对象的大小，说明客户端上传了对象的全部数据
		if current == stream.Size {
			// 调用Flush方法将剩余数据写入临时对象
			stream.Flush()
			// 调用NewRSResumableGetStream生成一个临时对象读取流getStream
			getStream, e := rs.NewRSResumableGetStream(stream.Servers, stream.Uuids, stream.Size)
			// 读取getStream中的数据并计算散列值
			hash := url.PathEscape(utils.CalculateHash(getStream))
			// 散列值不一致，客户端上传数据有错误
			if hash != stream.Hash {
				// 接口服务删除临时对象
				stream.Commit(false)
				log.Println("resumable put done but hash mismatch")
				// 返回403 Forbidden
				w.WriteHeader(http.StatusForbidden)
				return
			}
			// 检查该散列值是否已经存在
			if locate.Exist(url.PathEscape(hash)) {
				// 存在则删除临时对象
				stream.Commit(false)
			} else {
				// 否则将对象转正
				stream.Commit(true)
			}
			// 调用es.AddVersion添加新版本
			e = es.AddVersion(stream.Name, stream.Hash, stream.Size)
			if e != nil {
				log.Println(e)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
	}
}
