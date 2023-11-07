package ServerLog

import (
	"log"
	"os"
	"strconv"
	"time"
)

// serverGlobalLogCache 全局日志缓存, 用于写入文件.
var serverGlobalLogCache []string

// serverGlobalLogFileNameCount 全局日志文件计数器, 用于命名日志文件.
var serverGlobalLogFileNameCount uint64

// StarLogFileProcessEvent 启动日志处理[AsyncEvent].
// 定时器: 日志缓存 => 日志文件.
// @param timer, start uint64, folder, name string
// @return void
func StarLogFileProcessEvent(Timer, Start uint64, Folder, Name string) {
	if Timer > 10 {
		StarLogFmt("INFO", "LogFile Process Loop Start.", "Server", "Core", "Low", "Log", "File")
		// create Timer channel.
		ticker := time.Tick(time.Duration(Timer) * time.Second)

		// 设置日志文件编号起始值.
		cacheMutex.Lock()
		serverGlobalLogFileNameCount = Start
		cacheMutex.Unlock()

		// log process Timer event loop.
		for {
			select {
			case <-ticker:
				go writeServerLogFile(Folder, Name)
			}
		}
	} else {
		StarLogFmt("CRITICAL", "EventLoop Start Timer > 10s", "Server", "Core", "Low", "Log", "File")
	}
}

// writeServerLogFile 将全局缓存日志写入文件.
func writeServerLogFile(folder, name string) {
	LogFileName := folder
	cacheMutex.Lock()
	{
		LogFileName += strconv.FormatUint(serverGlobalLogFileNameCount, 10)
		serverGlobalLogFileNameCount++
	}
	cacheMutex.Unlock()
	LogFileName += name + ".log"
	StarLogFmt("INFO", "Write FileName: "+LogFileName, "Server", "Core", "Low", "Log", "FIle")

	// Open(create) log file.
	file, err := os.OpenFile(LogFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		StarLogFmt("ERROR", "File Open: "+err.Error(), "Server", "Core", "Low", "Log", "FIle")
	}

	// Func end => Close log file.
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			StarLogFmt("ERROR", "File Close: "+err.Error(), "Server", "Core", "Low", "Log", "FIle")
		}
	}(file)

	// Create Logger.
	Logger := log.New(file, "", log.Ldate|log.Ltime)

	cacheMutex.Lock()
	{
		// String array => Log file.
		for _, LogEntry := range serverGlobalLogCache {
			Logger.Print(LogEntry)
		}
		// Clear cache array data.
		serverGlobalLogCache = serverGlobalLogCache[:0]
	}
	cacheMutex.Unlock()
}
