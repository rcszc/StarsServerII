package ServerPanel

import (
	"ServerSTARS-II/ServerFrameworkCore/ServerLayerLow/ServerLog"
	"net"
	"sync"
	"time"
)

// ServerPanel UDP => Panel.

var stateSendString string
var stateLock sync.Mutex

// PanelSendUdpEvent 启动服务器面板(UDP)[事件循环]
// @param Str string
// @return void
func PanelSendUdpEvent(SendAddress string, Timer uint64) {
	if Timer >= 200 {
		TempStr := "PanelUDP Process Loop Start: " + SendAddress
		ServerLog.StarLogFmt("INFO", TempStr, "Server", "Core", "Panel", "Udp")
		// create timer.
		ticker := time.Tick(time.Duration(Timer) * time.Millisecond)

		// panel state process Timer event loop.
		for {
			select {
			case <-ticker:

				var TempStr []byte
				// string => bytes array.
				stateLock.Lock()
				{
					SendStr := panelStateProcess() + stateSendString
					TempStr = []byte(SendStr)
				}
				stateLock.Unlock()
				// 发送 UDP 数据包.
				go sendDataUDP(SendAddress, TempStr)
			}
		}
	} else {
		ServerLog.StarLogFmt("CRITICAL", "EventLoop Start Timer > 200ms", "Server", "Core", "Panel", "Udp")
	}
}

// PanelStateWrite 向写入面板状态 UDP 写入额外数据.
// @param Str string
// @return void
func PanelStateWrite(Str string) {
	stateLock.Lock()
	stateSendString = Str
	stateLock.Unlock()
}

// sendDataUDP 发送 UDP 数据包.
func sendDataUDP(address string, dataBytes []byte) {
	// 解析 UDP 地址.
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		ServerLog.StarLogFmt("WARNING", "Analysis UDP Address:"+err.Error(), "Server", "Core", "Panel", "Udp")
		return
	}

	// 建立 UDP 连接.
	connect, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return
	}
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {
			ServerLog.StarLogFmt("WARNING", "Create UDP:"+err.Error(), "Server", "Core", "Panel", "Udp")
		}
	}(connect)

	_, err = connect.Write(dataBytes)
	if err != nil {
		ServerLog.StarLogFmt("WARNING", "SendState UDP:"+err.Error(), "Server", "Core", "Panel", "Udp")
		return
	}
}
