package redis

// 定义一些常见的状态码

type Code int64

const (
	CodeSuccess Code = 1000 + iota
	CodeErrorRestore
	CodeMonitorAddrError
	CodeMonitorSentinelError
	CodeMonitorClusterError
)

var codeMsgMap = map[Code]string{
	CodeSuccess:              "success",
	CodeErrorRestore:         "故障恢复",
	CodeMonitorAddrError:     "redis 实例状态异常",
	CodeMonitorSentinelError: "redis 哨兵状态异常",
	CodeMonitorClusterError:  "redis 集群状态异常",
}

func (c Code) Message() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeSuccess]
	}
	return msg
}
