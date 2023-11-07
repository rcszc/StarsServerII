package ServerPanel

import (
	"ServerSTARS-II/ServerFrameworkCore/ServerLayerLow/ServerMonitor"
	"strconv"
)

// 数据格式[v20230821]:
// 分隔: " "(0x20), Float64
// [CPU总线程][CPU使用率][系统总内存][总内存使用率][服务器使用内存][硬盘总容量][硬盘总使用率][读取次数][写入次数][数据库QPS]

// panelStateProcess 服务器状态收集处理.
// data => string =UDP=> panel
// ServerMonitor => ServerPanel.
func panelStateProcess() string {
	var ReturnState string

	CpuState := ServerMonitor.MonitorGetCPU(false)
	MemState := ServerMonitor.MonitorGetMemory(false)
	DiskState := ServerMonitor.MonitorGetDisk(false)

	ReturnState += strconv.FormatFloat(float64(CpuState.CpuThreads), 'f', 3, 64) + " "
	ReturnState += strconv.FormatFloat(CpuState.CpuUsage, 'f', 3, 64) + " "

	ReturnState += strconv.FormatFloat(MemState.Capacity, 'f', 3, 64) + " "
	ReturnState += strconv.FormatFloat(MemState.Usage, 'f', 3, 64) + " "
	ReturnState += strconv.FormatFloat(MemState.System, 'f', 3, 64) + " "

	var DiCapacity, DiUsageSize float64
	var CountRead, CountWrite float64

	for _, disk := range DiskState {

		DiCapacity += disk.Capacity
		DiUsageSize += disk.Capacity * disk.Usage

		CountRead += float64(disk.IoCountRead)
		CountWrite += float64(disk.IoCountWrite)
	}
	ReturnState += strconv.FormatFloat(DiCapacity, 'f', 3, 64) + " "
	ReturnState += strconv.FormatFloat(DiUsageSize/DiCapacity, 'f', 3, 64) + " "

	ReturnState += strconv.FormatFloat(CountRead, 'f', 3, 64) + " "
	ReturnState += strconv.FormatFloat(CountWrite, 'f', 3, 64) + " "

	return ReturnState
}
