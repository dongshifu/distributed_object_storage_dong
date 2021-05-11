package es

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// 结构与ES映射中定义的objects类型属性一一对应
type Metadata struct {
	Name    string
	Version int
	Size    int64
	Hash    string
}

type hit struct {
	Source Metadata `json:"_source"`
}

type searchResult struct {
	Hits struct {
		Total int
		Hits  []hit
	}
}

// 根据对象的名字和版本号来获取对象的元数据
func getMetadata(name string, versionId int) (meta Metadata, e error) {
	// ES服务器地址来自环境变量ES_SERVER,索引是metadata,类型是objects
	// 文档的id由对象的名字和版本号拼接而成。
	url := fmt.Sprintf("http://%s/metadata/objects/%s_%d/_source",
		os.Getenv("ES_SERVER"), name, versionId)
	// GET到URL中的对象的元数据，免除耗时的搜索操作
	r, e := http.Get(url)
	if e != nil {
		return
	}
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("fail to get %s_%d: %d", name, versionId, r.StatusCode)
		return
	}
	// 读出数据
	result, _ := ioutil.ReadAll(r.Body)
	// ES返回的结果经过JSON解码后被es,SearchLatestVersion函数使用
	json.Unmarshal(result, &meta)
	return
}

func SearchLatestVersion(name string) (meta Metadata, e error) {
	// 调用ES搜索API.在URL中指定对象的名字，版本号以降序排列且只返回第一个结果。
	url := fmt.Sprintf("http://%s/metadata/_search?q=name:%s&size=1&sort=version:desc",
		os.Getenv("ES_SERVER"), url.PathEscape(name))
	fmt.Println(url)
	r, e := http.Get(url)
	if e != nil {
		return
	}
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("fail to search latest metadata: %d", r.StatusCode)
		return
	}
	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	// ES返回的结果通过JSON解码到一个searchResult结构体中
	// searchResult和ES搜索API返回的结构体保持一致
	// 方便读取搜索到的元数据并赋值给meta返回。
	json.Unmarshal(result, &sr)
	if len(sr.Hits.Hits) != 0 {
		meta = sr.Hits.Hits[0].Source
	}
	// ES的返回结果长度为0，说明没有搜到相对应的元数据，直接返回
	// 此时meta中各属性都为初始值：字符串为空，整型为0
	return
}

func GetMetadata(name string, version int) (Metadata, error) {
	// 当version为0的时候，调用SearchLatestVersion获取当前最新的版本
	if version == 0 {
		return SearchLatestVersion(name)
	}
	return getMetadata(name, version)
}

// 用于向ES服务上传一个新的元数据，输入的4个参数对应元数据的4个属性
func PutMetadata(name string, version int, size int64, hash string) error {
	// 生成ES文档，一个ES的文档相当于数据库的一条记录
	doc := fmt.Sprintf(`{"name":"%s","version":%d,"size":%d,"hash":"%s"}`,
		name, version, size, hash)
	client := http.Client{}
	// 使用op_type=create参数，如果同时又多个客户端上传同一个数据，结果会发生冲突
	// 只有第一个文档被成功创建，之后的PUT请求，ES会返回409Conflict
	url := fmt.Sprintf("http://%s/metadata/objects/%s_%d?op_type=create",
		os.Getenv("ES_SERVER"), name, version)
	// 用PUT方法将文档上传到metadata索引的objects类型
	request, _ := http.NewRequest("PUT", url, strings.NewReader(doc))
	r, e := client.Do(request)
	if e != nil {
		return e
	}
	// 如果为409Conflict，函数让版本号加1并递归调用自身继续上传
	if r.StatusCode == http.StatusConflict {
		return PutMetadata(name, version+1, size, hash)
	}
	if r.StatusCode != http.StatusCreated {
		result, _ := ioutil.ReadAll(r.Body)
		return fmt.Errorf("fail to put metadata: %d %s", r.StatusCode, string(result))
	}
	return nil
}

func AddVersion(name, hash string, size int64) error {
	// 获取对象最新的版本
	version, e := SearchLatestVersion(name)
	if e != nil {
		return e
	}
	// 在版本号上加1调用PutMetadata
	return PutMetadata(name, version.Version+1, size, hash)
}

// 用于搜索某个对象或所有对象的全部版本
func SearchAllVersions(name string, from, size int) ([]Metadata, error) {
	// name表示对象的名字，如果name不为空字符粗则搜索指定对象的所有版本
	// 否则搜索所有对象的全部版本
	// from和size指定分页的显示结果
	// 搜索的结果按照对象的名字和版本号排序
	url := fmt.Sprintf("http://%s/metadata/_search?sort=name,version&from=%d&size=%d",
		os.Getenv("ES_SERVER"), from, size)
	if name != "" {
		url += "&q=name:" + name
	}
	r, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	metas := make([]Metadata, 0)
	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	for i := range sr.Hits.Hits {
		metas = append(metas, sr.Hits.Hits[i].Source)
	}
	return metas, nil
}

// 根据对象的名字name和版本号version删除相应的对象元数据
func DelMetadata(name string, version int) {
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata/objects/%s_%d",
		os.Getenv("ES_SERVER"), name, version)
	request, _ := http.NewRequest("DELETE", url, nil)
	client.Do(request)
}

type Bucket struct {
	Key         string // 对象的名字
	Doc_count   int    //该对象目前有多少个版本
	Min_version struct {
		Value float32 // 当前最小的版本号
	}
}

type aggregateResult struct {
	Aggregations struct {
		Group_by_name struct {
			Buckets []Bucket
		}
	}
}

// 输入min_doc_count用于指示需要搜索对象最小版本数量
func SearchVersionStatus(min_doc_count int) ([]Bucket, error) {
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata/_search", os.Getenv("ES_SERVER"))
	// 使用ElasticSearch的aggregation search API搜索元数据
	// 以对象的名字分组，搜索版本数量大于等于min_doc_count的对象并返回
	body := fmt.Sprintf(`
        {
          "size": 0,
          "aggs": {
            "group_by_name": {
              "terms": {
                "field": "name",
                "min_doc_count": %d
              },
              "aggs": {
                "min_version": {
                  "min": {
                    "field": "version"
                  }
                }
              }
            }
          }
        }`, min_doc_count)
	request, _ := http.NewRequest("GET", url, strings.NewReader(body))
	r, e := client.Do(request)
	if e != nil {
		return nil, e
	}
	b, _ := ioutil.ReadAll(r.Body)
	var ar aggregateResult
	json.Unmarshal(b, &ar)
	return ar.Aggregations.Group_by_name.Buckets, nil
}

// 通过ES的search API搜索所有对象元数据中hash属性等于散列值的文档
func HasHash(hash string) (bool, error) {
	url := fmt.Sprintf("http://%s/metadata/_search?q=hash:%s&size=0", os.Getenv("ES_SERVER"), hash)
	r, e := http.Get(url)
	if e != nil {
		return false, e
	}
	b, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(b, &sr)
	// 如果满足条件的文档数量不为0，说明还存在对该散列值的引用，函数返回true,否则返回false
	return sr.Hits.Total != 0, nil
}

// 输入对象的散列值hash,通过ES的seach API查询元数据属性中hash等于该散列值的文档的size属性
func SearchHashSize(hash string) (size int64, e error) {
	url := fmt.Sprintf("http://%s/metadata/_search?q=hash:%s&size=1",
		os.Getenv("ES_SERVER"), hash)
	r, e := http.Get(url)
	if e != nil {
		return
	}
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("fail to search hash size: %d", r.StatusCode)
		return
	}
	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	if len(sr.Hits.Hits) != 0 {
		size = sr.Hits.Hits[0].Source.Size
	}
	// 返回size
	return
}
