package ServerTcpip

import (
	"ServerSTARS-II/ServerFrameworkCore/ServerLayerLow/ServerLog"
	"net/http"
	"strconv"
	"sync"
)

type WebSSL struct {
	FilePem string
	FileKey string
}

// TcpipHttpsServer 启动 HTTPS 服务器.
// @param Port uint16, SSL WebSSL, Wg *sync.WaitGroup
// @return void
func TcpipHttpsServer(Port uint16, SSL WebSSL, Wg *sync.WaitGroup) {
	Wg.Add(1)
	defer Wg.Done()
	ServerLog.StarLogFmt("INFO", "Start Server(HTTPS).", "Server", "Core", "Tcpip", "https")

	// 注册处理句柄(N).
	// 因为已经在 http 注册处理函数, 所以 https 无需注册.

	TempStr := ":" + strconv.FormatUint(uint64(Port), 10)
	// 启动 HTTPS 服务器.
	err := http.ListenAndServeTLS(TempStr, SSL.FilePem, SSL.FileKey, nil)
	if err != nil {
		ServerLog.StarLogFmt("ERROR", "Start HTTPS server(FAIL): "+err.Error(), "Server", "Core", "Tcpip", "https")
		return
	}
}
