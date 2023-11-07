package ServerLog

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var countMutex sync.Mutex
var cacheMutex sync.Mutex

var logModeNumberInfo uint64     // [mtx]
var logModeNumberWarning uint64  // [mtx]
var logModeNumberError uint64    // [mtx]
var logModeNumberCritical uint64 // [mtx]

const ConsColorDarkRed = "\033[31;2m"
const ConsColorRed = "\033[31m"
const ConsColorYellow = "\033[33m"
const ConsColorGray = "\033[90m"
const ConsColorReset = "\033[0m"

var consPrintColor string

// StarLogFmt 格式化标签输出日志.
// Mode: INFO(信息), WARNING(警告), ERROR(错误), CRITICAL(严重)
// @param mode, content_str, tag... string
// @return void
func StarLogFmt(Mode, Str string, Nums ...string) {
	var ResultStrLog = "[" + Mode + "]:"

	systemLogCount(Mode)

	for _, num := range Nums {
		ResultStrLog += "[" + num + "]:"
	}
	if len(Nums) > 0 {
		ResultStrLog += " " + Str
	} else {
		ResultStrLog += Str
	}

	CurrentTime := time.Now()
	// 使用神奇的 2006-01-02 15:04:05 格式化时间戳.
	OutCurrentTime := CurrentTime.Format("2006-01-02 15:04:05")
	ReturnLog := "[" + OutCurrentTime + "]" + ResultStrLog + "\n"

	cacheMutex.Lock()
	// Cache push string.
	serverGlobalLogCache = append(serverGlobalLogCache, ReturnLog)
	cacheMutex.Unlock()
	// Level color print.
	fmt.Printf(consPrintColor + ReturnLog + ConsColorReset)
}

// StarLogCount 输出日志统计信息.
// @param void
// @return string
func StarLogCount() string {
	var ResultCount string
	countMutex.Lock()
	{
		ResultCount += "[LogStatistics]:"
		ResultCount += " Info: " + strconv.FormatUint(logModeNumberInfo, 8)
		ResultCount += " Warning: " + strconv.FormatUint(logModeNumberWarning, 8)
		ResultCount += " Error: " + strconv.FormatUint(logModeNumberError, 8)
		ResultCount += " Critical: " + strconv.FormatUint(logModeNumberCritical, 8)
	}
	countMutex.Unlock()
	return ResultCount
}

func systemLogCount(mode string) {
	countMutex.Lock()
	{
		if mode == "INFO" {
			logModeNumberInfo++
			consPrintColor = ConsColorGray
		}
		if mode == "WARNING" {
			logModeNumberWarning++
			consPrintColor = ConsColorYellow
		}
		if mode == "ERROR" {
			logModeNumberError++
			consPrintColor = ConsColorRed
		}
		if mode == "CRITICAL" {
			logModeNumberCritical++
			consPrintColor = ConsColorDarkRed
		}
	}
	countMutex.Unlock()
}
