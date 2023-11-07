package ServerMonitor

import (
	"ServerSTARS-II/ServerFrameworkCore/ServerLayerLow/ServerLog"
	"github.com/shirou/gopsutil/cpu"
	"runtime"
	"strconv"
	"time"
)

type CpuInfo struct {
	CpuUsage   float64 // CPU使用率
	CpuThreads uint32  // CPU线程数
}

func MonitorGetCPU(PrintLog bool) (RetCpu CpuInfo) {
	SystemCpu, err := cpu.Info()
	if err != nil {
		ServerLog.StarLogFmt("ERROR", "SystemCpu Get Info: "+err.Error(), "Server", "Core", "Monitor", "CPU")
		return
	}

	var ReturnCpuInfo CpuInfo

	var CpuNumber, CpuCores, CpuLogicCores uint32
	var CpuFrequency float64 // AVG
	var CpuModel string

	// 获取 CPU 物理核心数.
	CoresTemp, err := cpu.Counts(false)
	if err != nil {
		ServerLog.StarLogFmt("ERROR", "Cpu Get Cores: "+err.Error(), "Server", "Core", "Monitor", "CPU")
		return
	}
	CpuCores = uint32(CoresTemp)

	// 获取 CPU 逻辑核心数.
	CoresTemp, err = cpu.Counts(true)
	if err != nil {
		ServerLog.StarLogFmt("ERROR", "Cpu Get LogicCores: "+err.Error(), "Server", "Core", "Monitor", "CPU")
		return
	}
	CpuLogicCores = uint32(CoresTemp)

	// 获取 CPU 数量,频率,型号.
	for _, info := range SystemCpu {
		CpuNumber++
		CpuFrequency += float64(info.Mhz) / 1000.0
		CpuModel += info.ModelName + " "
	}
	CpuFrequency /= float64(CpuNumber)

	// 获取 CPU 使用率.
	var CpuUsage = getCpuUsage()

	if PrintLog == true {
		ServerLog.StarLogFmt("INFO", "[CPU]: CPU: "+strconv.FormatUint(uint64(CpuNumber), 10)+", Model: "+CpuModel)

		TempStr := strconv.FormatUint(uint64(CpuCores), 10) + ","
		TempStr += " Logic: " + strconv.FormatUint(uint64(CpuLogicCores), 10) + ","
		TempStr += " AvgFreq: " + strconv.FormatFloat(CpuFrequency, 'f', 2, 64) + " GHz,"
		TempStr += " Usage: " + strconv.FormatFloat(CpuUsage, 'f', 2, 64) + " %%"

		ServerLog.StarLogFmt("INFO", "[CPU]: Cores: "+TempStr)
	}
	ReturnCpuInfo.CpuUsage = CpuUsage
	ReturnCpuInfo.CpuThreads = CpuLogicCores
	return ReturnCpuInfo
}

// getCpuUsage 获取CPU Cores平均使用率.
// @param void
// @return usage float64 (N%)
func getCpuUsage() (usage float64) {
	NumberCpu := runtime.NumCPU()

	// 获取每个 CPU 的利用率 (采样500ms).
	Percentages, err := cpu.Percent(500*time.Millisecond, false)
	if err != nil {
		ServerLog.StarLogFmt("ERROR", "SystemCpu Get SrcUsage: "+err.Error(), "Server", "Core", "Monitor", "CPU")
		return
	}

	// 平均(Avg)占用率.
	TotalUsage := 0.0
	for _, pct := range Percentages {
		TotalUsage += pct
	}
	return TotalUsage / float64(NumberCpu) * 100
}
