package SystemCoreFunction

type HttpDataIO struct {
	RequestURL  string // 请求信息(URL).
	RequestBody []byte // 请求正文.
	RequestType string // 请求文件类型.
	ReturnData  string // 返回(响应)数据.
}

// BusinessLogicProcessGET Server Core => Logic.
// 请求路径未命中缓存时 => 调用业务逻辑函数进行处理.
func BusinessLogicProcessGET(Process *HttpDataIO) {

}

// BusinessLogicProcessPOST Server Core => Logic.
func BusinessLogicProcessPOST(Process *HttpDataIO) {
	switch Process.RequestType {
	case "application/json":
	}
}
