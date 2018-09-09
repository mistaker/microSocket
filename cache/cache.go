package cache

import (
	"log"
	"microSocket/util"
	"sync"
	"time"
)

var (
	fileCache   sync.Map
	memoryCache sync.Map
	fileLock    sync.Mutex
)

//缓存 数据
func SaveCache(key, value string, isFile bool) bool {
	if isFile {
		_, ok := fileCache.LoadOrStore(key, value)
		return ok
	} else {
		_, ok := memoryCache.LoadOrStore(key, value)
		return ok
	}

}

//获取数据
func GetCache(key string, isFile bool) string {
	if isFile {
		v, ok := fileCache.Load(key)
		if ok {
			return v.(string)
		}
	} else {
		v, ok := memoryCache.Load(key)
		if ok {
			return v.(string)
		}
	}
	return "'"
}

//删除数据
func DeleteCache(key string, isFile bool) {
	if isFile {
		fileCache.Delete(key)
	} else {
		memoryCache.Delete(key)
	}
}

//每隔一段时间缓存一次到文件里面
func saveFile() {
	saveMap := make(map[string]string)

	fileCache.Range(func(key, value interface{}) bool {
		saveMap[key.(string)] = value.(string)
		return true
	})

	saveStr := util.Map2String(saveMap)

	fileLock.Lock()
	defer fileLock.Unlock()
	log.Println(saveStr)
	util.WriteFile("./tmp.txt", saveStr)
}

//在此模块一开始运行的时候加载缓存文件
func loadFile() {
	fileLock.Lock()
	defer fileLock.Unlock()

	fileStr, _ := util.ReadFile("./tmp.txt")
	loadMap := util.String2Map(fileStr)

	for i, v := range loadMap {
		fileCache.LoadOrStore(i, v)
	}
}

func init() {
	loadFile()

	go func() {
		for {
			time.Sleep(time.Second)
			saveFile()
		}
	}()
}
