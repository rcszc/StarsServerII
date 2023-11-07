package ServerMonitor

import (
	"ServerSTARS-II/ServerFrameworkCore/ServerLayerLow/ServerLog"
	"github.com/shirou/gopsutil/mem"
	"runtime"
	"strconv"
)

type MemoryInfo struct {
	Capacity float64 // 内存容量(MiB)
	Usage    float64 // 内存容量占用率(%)
	System   float64 // 当前系统分配内存(MiB)
}

// MonitorGetMemory 监控内存[使用外部库, RunTime]
// @param PrintLog bool
// @return RetMemory MemoryInfo
func MonitorGetMemory(PrintLog bool) (RetMemory MemoryInfo) {
	MemInfo, err := mem.VirtualMemory()
	if err != nil {
		ServerLog.StarLogFmt("ERROR", "SystemMemory Get Info: "+err.Error(), "Server", "Core", "Monitor", "Mem")
		return
	}

	var ReturnMemoryInfo MemoryInfo

	// runtime get system memory usage
	var MemStats runtime.MemStats
	runtime.ReadMemStats(&MemStats)
	var SystemMemory = float64(MemStats.Sys) / (1024.0 * 1024.0)

	var TotalMemory = float64(MemInfo.Total) / (1024.0 * 1024.0 * 1024.0)
	var MemoryUsage = float64(MemInfo.Used) / (1024.0 * 1024.0 * 1024.0) / TotalMemory * 100

	if PrintLog == true {
		var PrintTemp = "[Memory]:"

		PrintTemp += " Capacity: " + strconv.FormatFloat(TotalMemory, 'f', 3, 64) + " GiB,"
		PrintTemp += " Usage: " + strconv.FormatFloat(MemoryUsage, 'f', 2, 64) + "%%,"
		PrintTemp += " System: " + strconv.FormatFloat(SystemMemory, 'f', 2, 64) + " MiB"

		ServerLog.StarLogFmt("INFO", PrintTemp)
	}

	ReturnMemoryInfo.Capacity = float64(MemInfo.Total) / (1024.0 * 1024.0)
	ReturnMemoryInfo.Usage = MemoryUsage
	ReturnMemoryInfo.System = SystemMemory

	return ReturnMemoryInfo
}
