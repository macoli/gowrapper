package redisfunc

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func CheckSentinelManager(addr, password string) {
	errMsg := Msg{
		Code:  CodeMonitorSentinelError,
		Title: CodeMonitorSentinelError.Message(),
		Time:  time.Now().Format("2006-01-02 15:04:05"),
	}

	// 创建连接
	rc, err := InitSentinelManageRedis(addr, password)
	if err != nil {
		errMsg.Err = err
		errMsg.Notice()
	}
	defer rc.Close()

	// 订阅哨兵节点的 channel
	pubSub := rc.Subscribe(ctx, sentinelChannels...)
	_, err = pubSub.Receive(ctx)
	if err != nil {
		tmp := fmt.Sprintf("获取 redis 哨兵节点 %s 的 channel 失败: %v\n", addr, err)

		errMsg.Err = errors.New(tmp)
		errMsg.Notice()
	}
	ch := pubSub.Channel()

	// 处理接收到的消息
	for message := range ch {
		//fmt.Println(msg.Channel, msg.Payload, "\r\n")
		// 处理每个 channel 接收到的消息
		switch message.Channel {
		case "+switch-master": // 发生了主从切换
			contentTMP := strings.Split(message.Payload, " ")
			tmp := fmt.Sprintf("master name: %s 发生了主从切换,master 从 %v:%v 切换为 %v:%v",
				contentTMP[0], contentTMP[1], contentTMP[2], contentTMP[3], contentTMP[4])

			errMsg.Err = errors.New(tmp)
			errMsg.Notice()
		default: // 其他哨兵事件通知
			msg := Msg{
				Code:    CodeSuccess,
				Title:   "redis 哨兵事件提醒",
				Time:    time.Now().Format("2006-01-02 15:04:05"),
				Content: fmt.Sprintf("%v %v", message.Channel, message.Payload),
			}
			msg.Print()
		}
	}

}

func CheckSentinel(sentinels []string, password, masterName string) {
	errMsg := Msg{
		Code:  CodeMonitorSentinelError,
		Title: CodeMonitorSentinelError.Message(),
		Time:  time.Now().Format("2006-01-02 15:04:05"),
	}

	// master name 对应 master 实例
	// 创建连接
	rm, err := InitSentinelMaster(sentinels, password, masterName)
	if err != nil {
		errMsg.Err = err
		errMsg.Notice()
	}
	defer rm.Close()

	// master name 对应 slave 实例
	// 创建连接
	rs, err := InitSentinelSlave(sentinels, password, masterName)
	if err != nil {
		errMsg.Err = err
		errMsg.Notice()
	}
	defer rs.Close()

}
