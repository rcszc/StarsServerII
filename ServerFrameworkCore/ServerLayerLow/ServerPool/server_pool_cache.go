package ServerPool

import (
	"ServerSTARS-II/ServerFrameworkCore/ServerLayerLow/ServerLog"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

// globalHashCache 全局缓存 HashMap.
var globalHashCache = make(map[string]string)

// globalHashCacheSize 全局缓存 size 计数.
var globalHashCacheSize uint64

// cacheMutex 全局缓存锁.
var cacheMutex sync.Mutex

// CachePushFile 向全局缓存添加文件.
// WARN: 每个文件大小不得超过 1024MiB (PSA)
// @param Key, FileName string
// @return void
func CachePushFile(Key, FileName string) {
	FileState, err := os.Stat(FileName)
	if err != nil {
		ServerLog.StarLogFmt("ERROR", "File State: "+err.Error(), "Server", "Core", "Low", "Pool", "Cache")
		return
	} else if FileState.Size() < (1024 * 1024 * 1024) {
		// File < 1GiB.
		// 文件路径 + 文件大小 + 文件扩展名.
		TempStr := strconv.FormatFloat(float64(FileState.Size())/(1024.0), 'f', 2, 64) + " KiB ["
		TempStr += filepath.Ext(FileName) + "] "
		ServerLog.StarLogFmt("INFO", "Load File: "+FileName+" "+TempStr, "Server", "Core", "Low", "Pool", "Cache")

		file, err := os.Open(FileName)
		if err != nil {
			ServerLog.StarLogFmt("ERROR", "File Open: "+err.Error(), "Server", "Core", "Low", "Pool", "Cache")
			return
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				ServerLog.StarLogFmt("ERROR", "File Close: "+err.Error(), "Server", "Core", "Low", "Pool", "Cache")
			}
		}(file)

		BinaryData, err := io.ReadAll(file)
		if err != nil {
			ServerLog.StarLogFmt("ERROR", "Read BinFile: "+err.Error(), "Server", "Core", "Low", "Pool", "Cache")
			return
		}

		// Binary => Base64(string) => HashMap. [废弃]
		// Mew: Binary => HashMap.
		pushMspString(Key, string(BinaryData))
	} else {
		ServerLog.StarLogFmt("ERROR", "Failed Load, FileSize > 1024MiB.", "Server", "Core", "Low", "Pool", "Cache")
	}
}

// CacheFind 查询 HashMap 缓存.
// @param Key string, Exists *bool
// @return string
func CacheFind(Key string, Exists *bool) string {
	var ReturnString string
	var ExistsValue bool
	cacheMutex.Lock()
	{
		ReturnString, ExistsValue = globalHashCache[Key]
	}
	cacheMutex.Unlock()
	// set find cache state.
	Exists = &ExistsValue
	return ReturnString
}

// CachePushString 向全局缓存添加 string.
// @param Key, StrData string
// @return void
func CachePushString(Key, StrData string) {
	TempStr := strconv.FormatFloat(float64(len(StrData)/(1024.0*1024.0)), 'f', 2, 64) + " MiB"
	ServerLog.StarLogFmt("INFO", "Load StringSize: "+TempStr, "Server", "Core", "Low", "Pool", "Cache")
	pushMspString(Key, StrData)
}

// CacheDeleteData 从全局缓存中删除数据.
// @param Key string
// @return void
func CacheDeleteData(Key string) {
	TempStr := strconv.FormatFloat(float64(deleteMspString(Key)/(1024.0*1024.0)), 'f', 2, 64) + " MiB"
	ServerLog.StarLogFmt("INFO", "Delete CacheSize: "+TempStr, "Server", "Core", "Low", "Pool", "Cache")
}

// CacheSizeLen 查询缓存大小(items)
// @param PrintLog bool
// @return void
func CacheSizeLen(PrintLog bool) uint32 {
	var ReturnSizeLen uint32
	cacheMutex.Lock()
	{
		ReturnSizeLen = uint32(len(globalHashCache))
	}
	cacheMutex.Unlock()
	// print switch.
	if PrintLog {
		TempStr := strconv.FormatUint(uint64(ReturnSizeLen), 10)
		ServerLog.StarLogFmt("INFO", "Cache Length: "+TempStr+" items", "Server", "Core", "Low", "Pool", "Cache")
	}
	return ReturnSizeLen
}

// pushMspString 全局缓存添加数据 push + count.
func pushMspString(key, data string) {
	cacheMutex.Lock()
	{
		globalHashCacheSize += uint64(len(data))
		globalHashCache[key] = data
	}
	cacheMutex.Unlock()
}

// deleteMspString 全局缓存删除数据 delete data + count.
func deleteMspString(key string) uint64 {
	var ReturnDelSize uint64
	cacheMutex.Lock()
	{
		ReturnDelSize = uint64(len(globalHashCache[key]))
		globalHashCacheSize -= ReturnDelSize
		delete(globalHashCache, key)
	}
	cacheMutex.Unlock()
	return ReturnDelSize
}
