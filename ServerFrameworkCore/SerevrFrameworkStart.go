package ServerFrameworkCore

import (
	"ServerSTARS-II/ServerFrameworkCore/ServerLayerLow/ServerDatabase"
	"ServerSTARS-II/ServerFrameworkCore/ServerLayerLow/ServerLog"
	"ServerSTARS-II/ServerFrameworkCore/ServerLayerLow/ServerMonitor"
	"ServerSTARS-II/ServerFrameworkCore/ServerLayerLow/ServerPool"
	"ServerSTARS-II/ServerFrameworkCore/ServerLayerLow/ServerTool"
	"ServerSTARS-II/ServerFrameworkCore/ServerLayerMiddle/ServerPanel"
	"ServerSTARS-II/ServerFrameworkCore/ServerLayerMiddle/ServerTcpip"
	"sync"
)

// [2023_08_14] RCSZ.
// Core外部依赖库:
// 服务器监控 [磁盘监控][内存监控][CPU监控]
// github.com/shirou/gopsutil/disk
// github.com/shirou/gopsutil/mem
// github.com/shirou/gopsutil/cpu
// 服务器数据库
// MySQL(MariaDB) 驱动 github.com/go-sql-driver/mysql

var MainThread sync.WaitGroup

func MainStartServer(Name string) {
	ServerLog.StarLogFmt("INFO", "\033[96mPomeloStar Server Start: "+Name, "Server", "Core")

	// 检查注册文件完整性.
	if ConfigFilesState("ServerSystem/FILES_REGISTRATION.json") == true {
		// HTTP & HTTPS 注册相同处理函数(同池).
		var HttpPcsPool *ServerPool.GoPool

		// 读取 Core 配置.
		TempStr, Find := RegFileHashMap["FileSTCCORE"].(string)
		if Find == true {
			CoreConfig := ServerTool.JsonFileLoad(TempStr)

			// 启动日志文件处理事件.
			go ServerLog.StarLogFileProcessEvent(
				uint64(CoreConfig["LogWriteFileTimer"].(float64)),
				uint64(CoreConfig["LogWriteFileStart"].(float64)),
				CoreConfig["LogWriteFileFolder"].(string),
				CoreConfig["LogWriteFileName"].(string),
			)

			// 启动面板状态处理事件.
			go ServerPanel.PanelSendUdpEvent(
				CoreConfig["StateSendAddress"].(string),
				uint64(CoreConfig["StateSendTimerMs"].(float64)),
			)

			// 加载网页静态资源.
			WebStaticLoadResources(CoreConfig["WebStaticResFolder"].(string))
		} else {
			ServerLog.StarLogFmt("WARNING", "InitKey Loss: 'FileSTCCORE'", "Server", "Core", "Start")
		}

		// Print system info state.
		ServerMonitor.MonitorGetCPU(true)
		ServerMonitor.MonitorGetMemory(true)
		ServerMonitor.MonitorGetDisk(true)
		ServerPool.CacheSizeLen(true)

		// 读取数据库配置(MariaDB).
		TempStr, Find = RegFileHashMap["FileSTCDB"].(string)
		if Find == true {
			DBConfig := ServerTool.JsonFileLoad(TempStr)
			// 连接到数据库.
			DataBase, err := ServerDatabase.MariaDBConnect(
				DBConfig["DBUserName"].(string),
				DBConfig["DBPassword"].(string),
				DBConfig["DBHost"].(string),
				DBConfig["DBPort"].(string),
				DBConfig["DataBaseName"].(string),
			)
			if err == nil {
				ServerDatabase.MariaDBProcess(DataBase)
			}
		} else {
			ServerLog.StarLogFmt("WARNING", "InitKey Loss: 'FileSTCDB'", "Server", "Core", "Start")
		}

		// 读取 http 服务器配置.
		TempStr, Find = RegFileHashMap["FileSTCHTTP"].(string)
		if Find == true {
			HttpConfig := ServerTool.JsonFileLoad(TempStr)

			HttpPcsPool = ServerPool.CreateGoPool(uint32(HttpConfig["ProcessPoolMax"].(float64)))
			// 启动 https 服务器(async).
			go ServerTcpip.TcpipHttpServer(
				HttpPcsPool,
				uint32(HttpConfig["ProcessPoolWarn"].(float64)),
				uint16(HttpConfig["ServerPort"].(float64)),
				&MainThread,
			)
		} else {
			ServerLog.StarLogFmt("WARNING", "InitKey Loss: 'FileSTCHTTP'", "Server", "Core", "Start")
		}

		// 读取 https 服务器配置.
		TempStr, Find = RegFileHashMap["FileSTCHTTPS"].(string)
		if Find == true {
			HttpsConfig := ServerTool.JsonFileLoad(TempStr)

			// 加载 SSL 证书文件路径.
			var TempFileName ServerTcpip.WebSSL
			TempFileName.FilePem = HttpsConfig["WebSSLFileNginxPem"].(string)
			TempFileName.FileKey = HttpsConfig["WebSSLFileNginxKey"].(string)
			//启动 https 服务器(async).
			go ServerTcpip.TcpipHttpsServer(uint16(HttpsConfig["ServerPort"].(float64)), TempFileName, &MainThread)
		} else {
			ServerLog.StarLogFmt("WARNING", "InitKey Loss: 'FileSTCHTTPS'", "Server", "Core", "Start")
		}
		// 阻塞主线程.
		MainThread.Wait()

		// Close End(Free).
		MainFreeStarServer(HttpPcsPool)
	} else {
		ServerLog.StarLogFmt("CRITICAL", "REG FILE...")
	}
}

func MainFreeStarServer(http *ServerPool.GoPool) {
	http.Close()
}
