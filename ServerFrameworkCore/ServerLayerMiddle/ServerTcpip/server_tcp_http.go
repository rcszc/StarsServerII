package ServerTcpip

import (
	"ServerSTARS-II/ServerFrameworkCore/ServerLayerLow/ServerLog"
	"ServerSTARS-II/ServerFrameworkCore/ServerLayerLow/ServerPool"
	"ServerSTARS-II/ServerFrameworkLogic/SystemCoreFunction"
	"io"
	"net/http"
	"strconv"
	"sync"
)

// RequestProcessPoolWarn 请求任务处理池 负载警告标志.
var RequestProcessPoolWarn bool
var ReqPoolFlagMutex sync.Mutex

// TcpipHttpServer 启动 HTTP 服务器.
// @param PcsPool *ServerPool.GoPool, PcsWarn uint32, Port uint16, Wg *sync.WaitGroup
// @return void
func TcpipHttpServer(PcsPool *ServerPool.GoPool, PcsWarn uint32, Port uint16, Wg *sync.WaitGroup) {
	Wg.Add(1)
	defer Wg.Done()
	// 启动 http 请求处理协程池.
	PcsPool.Start()
	ServerLog.StarLogFmt("INFO", "Start Server(HTTP).", "Server", "Core", "Tcpip", "http")

	// 注册处理句柄(函数).
	// Request => Process Pool.
	http.HandleFunc("/", func(wrt http.ResponseWriter, req *http.Request) {
		// 将请求处理函数.
		httpRequestProcess(PcsPool, wrt, req)

		// Task Process WarnFlag.
		ReqPoolFlagMutex.Lock()
		{
			if PcsPool.TasksNumber() > PcsWarn {
				RequestProcessPoolWarn = true
			} else {
				RequestProcessPoolWarn = false
			}
		}
		ReqPoolFlagMutex.Unlock()
	})

	TempStr := ":" + strconv.FormatUint(uint64(Port), 10)
	// 启动 HTTP 服务器.
	err := http.ListenAndServe(TempStr, nil)
	if err != nil {
		ServerLog.StarLogFmt("ERROR", "Start HTTP server(FAIL): "+err.Error(), "Server", "Core", "Tcpip", "http")
		return
	}
}

// TcpipFindProcessWarn 查询任务处理池是否达到警告线.
// @param void
// @return bool
func TcpipFindProcessWarn() bool {
	var ReturnWarnFlag bool
	ReqPoolFlagMutex.Lock()
	{
		ReturnWarnFlag = RequestProcessPoolWarn
	}
	ReqPoolFlagMutex.Unlock()
	return ReturnWarnFlag
}

// httpRequestProcess 包含 HTTPS 解码后的请求.
// @param w http.ResponseWriter, r *http.Request
// @return void
func httpRequestProcess(pool *ServerPool.GoPool, w http.ResponseWriter, r *http.Request) {
	StringTempURL := r.URL.Path
	// request body => []byte
	BodyText, err := io.ReadAll(r.Body)
	if err != nil {
		ServerLog.StarLogFmt("ERROR", "Read ReqBody: "+err.Error(), "Server", "Core", "Tcpip", "http")
	}
	// Server Core => 业务逻辑.
	var ProcessTask SystemCoreFunction.HttpDataIO

	// "GET" request.
	if r.Method == http.MethodGet {
		ServerLog.StarLogFmt("INFO", "'GET' Req: "+r.RemoteAddr+" "+StringTempURL, "Server", "Core", "Tcpip", "http")

		var ReturnStaticTemp string
		var FindCacheState bool

		// 无路径直接请求主网页(不带参).
		if (StringTempURL == "/") && (len(BodyText) == 0) {
			// Find CacheMap.
			ReturnStaticTemp = ServerPool.CacheFind("/index.html", &FindCacheState)
		} else {
			ReturnStaticTemp = ServerPool.CacheFind(StringTempURL, &FindCacheState)
		}

		// 缓存无静态资源 & 带请求参数 => 调用业务处理逻辑.
		if FindCacheState || (len(BodyText) > 0) {
			ServerLog.StarLogFmt("INFO", "'GET' Missed Cache => BL.", "Server", "Core", "Tcpip", "http")

			// r.URL.Path => URLTemp => ProcessTask
			ProcessTask.RequestURL = StringTempURL
			ProcessTask.RequestBody = BodyText
			ProcessTask.RequestType = r.Header.Get("Content-Type")
			// => TaskPointer
			SystemCoreFunction.BusinessLogicProcessGET(&ProcessTask)
			ReturnStaticTemp = ProcessTask.ReturnData
		}

		// 设置 Content 响应头.
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.Header().Set("Content-Length", strconv.Itoa(len(ReturnStaticTemp)))

		// 返回(响应)数据.
		_, err := w.Write([]byte(ReturnStaticTemp))
		if err != nil {
			ServerLog.StarLogFmt("ERROR", "'GET' ReqRet: "+err.Error(), "Server", "Core", "Tcpip", "http")
			return
		}
	}

	// "POST" request.
	if r.Method == http.MethodPost {

		ProcessTask.RequestURL = StringTempURL
		ProcessTask.RequestBody = BodyText
		ProcessTask.RequestType = r.Header.Get("Content-Type")

		TempStr := " Size: " + strconv.FormatUint(uint64(len(BodyText)), 10) + " Bytes"
		ServerLog.StarLogFmt("INFO", "'POST' Req: "+r.RemoteAddr+TempStr, "Server", "Core", "Tcpip", "http")

		SystemCoreFunction.BusinessLogicProcessPOST(&ProcessTask)

		_, err = w.Write([]byte(ProcessTask.ReturnData))
		if err != nil {
			ServerLog.StarLogFmt("ERROR", "'POST' ReqRet: "+err.Error(), "Server", "Core", "Tcpip", "http")
			return
		}
	}
}
