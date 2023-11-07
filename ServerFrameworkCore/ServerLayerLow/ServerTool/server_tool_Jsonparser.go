package ServerTool

import (
	"ServerSTARS-II/ServerFrameworkCore/ServerLayerLow/ServerLog"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// JsonFileLoad 读取 JSON 配置文件到哈希表.
// @param FileName string
// @return map[string]interface{}, error
func JsonFileLoad(FileName string) map[string]interface{} {
	ServerLog.StarLogFmt("INFO", "Read JSON Config: "+FileName, "Server", "Core", "Config")
	var ReturnHashMap map[string]interface{}
	// 打开 JSON 文件.
	file, err := os.Open(FileName)
	if err != nil {
		ServerLog.StarLogFmt("ERROR", "File Open: "+err.Error(), "Server", "Core", "Config")
		return nil
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			ServerLog.StarLogFmt("ERROR", "File Close: "+err.Error(), "Server", "Core", "Config")
		}
	}(file)

	// 读取 JSON 文件.
	JsonData, err := io.ReadAll(file)
	if err != nil {
		ServerLog.StarLogFmt("ERROR", "Read JsonFile: "+err.Error(), "Server", "Core", "Config")
		return nil
	}

	// 解码 JSON 数据 => HashMap.
	err = json.Unmarshal(JsonData, &ReturnHashMap)
	if err != nil {
		ServerLog.StarLogFmt("ERROR", "Decode JsonFile: "+err.Error(), "Server", "Core", "Config")
		return nil
	}
	return ReturnHashMap
}

// JsonFmtHashMapString 格式化 map[string]interface{} 为 map[string]string
// 为兼容 ServerConfig.JsonFileLoad
func JsonFmtHashMapString(InputMap map[string]interface{}) map[string]string {
	// Create map[string]string
	var ReturnHashMap = make(map[string]string)

	// 遍历原始 Map 并转换.
	for key, value := range InputMap {
		strValue, StrOK := value.(string)
		if StrOK {
			ReturnHashMap[key] = strValue
		} else {
			// interface{} => String
			ReturnHashMap[key] = fmt.Sprintf("%v", value)
		}
	}
	return ReturnHashMap
}
