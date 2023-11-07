package ServerTcpip

import (
	"ServerSTARS-II/ServerFrameworkCore/ServerLayerLow/ServerLog"
	"encoding/json"
)

// HTTP 请求处理工具.

// ToolBodyToJsonHashMap 请求正文数据以 Json 解析.
func ToolBodyToJsonHashMap(HttpBody []byte) map[string]interface{} {
	var ReturnHashMap map[string]interface{}
	// Decode Json data => HashMap.
	err := json.Unmarshal(HttpBody, &ReturnHashMap)
	if err != nil {
		ServerLog.StarLogFmt("ERROR", "Decode JSON: "+err.Error(), "Server", "Core", "Tcpip", "tool")
		return nil
	}
	return ReturnHashMap
}
