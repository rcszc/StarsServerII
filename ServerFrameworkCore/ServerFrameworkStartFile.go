package ServerFrameworkCore

import (
	"ServerSTARS-II/ServerFrameworkCore/ServerLayerLow/ServerLog"
	"ServerSTARS-II/ServerFrameworkCore/ServerLayerLow/ServerTool"
	"os"
	"strconv"
)

var configFilesList = []string{"FileSTCCORE", "FileSTCDBSQL", "FileSTCDB", "FileSTCHTTP", "FileSTCHTTPS"}

// RegFileHashMap 服务器系统配置文件注册表.
var RegFileHashMap map[string]interface{}

// ConfigFilesState 检查系统配置文件完整性.
// @param REGFILE string
// @return bool
func ConfigFilesState(REGFILE string) bool {
	ServerLog.StarLogFmt("INFO", "REG Files Find...", "Server", "Core", "StartFile")
	RegFileHashMap = ServerTool.JsonFileLoad(REGFILE)

	var ReturnFilesStat = true
	var FilesExists uint32

	for _, FileName := range configFilesList {
		// 查询 HashMap.
		FileNameTemp, _ := RegFileHashMap[FileName].(string)

		_, err := os.Stat(FileNameTemp)
		if err == nil {
			FilesExists++
		} else if os.IsNotExist(err) {
			ReturnFilesStat = false
		} else {
			ServerLog.StarLogFmt("ERROR", "FileStat Get: "+err.Error(), "Server", "Core", "StartFile")
		}
	}
	TempStr := strconv.FormatUint(uint64(FilesExists), 10) + " items"
	ServerLog.StarLogFmt("INFO", "REG Files Exists: "+TempStr, "Server", "Core", "StartFile")
	return ReturnFilesStat
}
