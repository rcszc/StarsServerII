package ServerMonitor

import (
	"ServerSTARS-II/ServerFrameworkCore/ServerLayerLow/ServerLog"
	"github.com/shirou/gopsutil/disk"
	"strconv"
)

type DiskInfo struct {
	Capacity     float64 // 分区容量(GiB)
	Usage        float64 // 分区容量使用率(%)
	IoCountRead  uint64  // 分区读出次数
	IoCountWrite uint64  // 分区写入次数
}

// MonitorGetDisk 监控磁盘[使用外部库]
// @param PrintLog bool
// @return RetDisk []DiskInfo
func MonitorGetDisk(PrintLog bool) (RetDisk []DiskInfo) {
	// 获取磁盘信息.
	Partitions, err := disk.Partitions(false)
	if err != nil {
		ServerLog.StarLogFmt("ERROR", "Disk Get SrcInfo: "+err.Error(), "Server", "Core", "Monitor", "Disk")
		return
	}

	var ReturnDiskInfo []DiskInfo

	// 遍历每个硬盘分区.
	for _, Partition := range Partitions {
		Usage, err := disk.Usage(Partition.Mountpoint)
		if err != nil {
			ServerLog.StarLogFmt("ERROR", "SystemDisk Get Info(s): "+err.Error(), "Server", "Core", "Monitor", "Disk")
			continue
		}

		// 获取磁盘 IO 信息.
		IoUsage, err := disk.IOCounters(Partition.Mountpoint)
		if err != nil {
			ServerLog.StarLogFmt("ERROR", "SystemDisk Get IOCount: "+err.Error(), "Server", "Core", "Monitor", "Disk")
			return
		}

		// 打印磁盘信息 src.Bytes => GBytes.
		if PrintLog == true {
			var PrintTemp = "[Disk]: " + Partition.Mountpoint
			var FloatTemp = float64(Usage.Total) / (1024.0 * 1024.0 * 1024.0)

			PrintTemp += " Capacity: " + strconv.FormatFloat(FloatTemp, 'f', 2, 64) + " GiB,"
			PrintTemp += " Usage: " + strconv.FormatFloat(Usage.UsedPercent, 'f', 2, 64) + "%%"
			PrintTemp += " Read: " + strconv.FormatUint(IoUsage[Partition.Mountpoint].ReadCount, 10)
			PrintTemp += " Write: " + strconv.FormatUint(IoUsage[Partition.Mountpoint].WriteCount, 10)

			ServerLog.StarLogFmt("INFO", PrintTemp)
		}

		var InfoTemp DiskInfo

		InfoTemp.Capacity = float64(Usage.Total) / (1024.0 * 1024.0 * 1024.0)
		InfoTemp.Usage = Usage.UsedPercent
		InfoTemp.IoCountRead = IoUsage[Partition.Mountpoint].ReadCount
		InfoTemp.IoCountWrite = IoUsage[Partition.Mountpoint].WriteCount

		ReturnDiskInfo = append(ReturnDiskInfo, InfoTemp)
	}
	return ReturnDiskInfo
}
