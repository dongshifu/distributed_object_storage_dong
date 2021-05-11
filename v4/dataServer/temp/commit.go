package temp

import (
	"net/url"
	"object_storage_dong/lib/utils"
	"object_storage_dong/v4/dataServer/locate"
	"os"
	"strconv"
	"strings"
)

func (t *tempInfo) hash() string {
	s := strings.Split(t.Name, ".")
	return s[0]
}

func (t *tempInfo) id() int {
	s := strings.Split(t.Name, ".")
	id, _ := strconv.Atoi(s[1])
	return id
}

func commitTempObject(datFile string, tempinfo *tempInfo) {
	// 读取临时对象的数据并计算散列值<hash of shard X>
	f, _ := os.Open(datFile)
	d := url.PathEscape(utils.CalculateHash(f))
	f.Close()
	os.Rename(datFile, os.Getenv("STORAGE_ROOT")+"/objects/"+tempinfo.Name+"."+d)
	// 调用locate.Add,以<hash>为键，分片的id为值添加进定位缓存
	locate.Add(tempinfo.hash(), tempinfo.id())
}
