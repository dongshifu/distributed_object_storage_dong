package main

import (
	"log"
	"object_storage_dong/lib/es"
)

const MIN_VERSION_COUNT = 5

func main() {
	// 调用es.SearchVersionStatus将元数据服务中所有版本数量大于等于6的对象搜索出来保存到Bucket结构体的数组buckets中
	buckets, e := es.SearchVersionStatus(MIN_VERSION_COUNT + 1)
	if e != nil {
		log.Println(e)
		return
	}
	// 遍历buckets
	for i := range buckets {
		bucket := buckets[i]
		// 在for循环中调用es.DelMetadata,从该对象当前最小的版本号开始一一删除
		// 直到最后剩下5个
		for v := 0; v < bucket.Doc_count-MIN_VERSION_COUNT; v++ {
			es.DelMetadata(bucket.Key, v+int(bucket.Min_version.Value))
		}
	}
}
