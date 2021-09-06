package redis

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/macoli/iwrapper/message"
)

var sentinelChannels = []string{
	"+reset-master",                     // 主服务器已被重置。
	"+slave",                            // 一个新的从服务器已经被 Sentinel 识别并关联。
	"+failover-state-reconf-slaves",     // 故障转移状态切换到了 reconf-slaves 状态。
	"+failover-detected",                // 另一个 Sentinel 开始了一次故障转移操作，或者一个从服务器转换成了主服务器。
	"+slave-reconf-sent",                // 领头（leader）的 Sentinel 向实例发送了 [SLAVEOF](/commands/slaveof.html) 命令，为实例设置新的主服务器。
	"+slave-reconf-inprog",              // 实例正在将自己设置为指定主服务器的从服务器，但相应的同步过程仍未完成。
	"+slave-reconf-done",                // 从服务器已经成功完成对新主服务器的同步。
	"-dup-sentinel",                     // 对给定主服务器进行监视的一个或多个 Sentinel 已经因为重复出现而被移除 —— 当 Sentinel 实例重启的时候，就会出现这种情况。
	"+sentinel",                         // 一个监视给定主服务器的新 Sentinel 已经被识别并添加。
	"+sdown",                            // 给定的实例现在处于主观下线状态。
	"-sdown",                            // 给定的实例已经不再处于主观下线状态。
	"+odown",                            // 给定的实例现在处于客观下线状态。
	"-odown",                            // 给定的实例已经不再处于客观下线状态。
	"+new-epoch",                        // 当前的纪元（epoch）已经被更新。
	"+try-failover",                     // 一个新的故障迁移操作正在执行中，等待被大多数 Sentinel 选中（waiting to be elected by the majority）。
	"+elected-leader",                   // 赢得指定纪元的选举，可以进行故障迁移操作了。
	"+failover-state-select-slave",      // 故障转移操作现在处于 select-slave 状态 —— Sentinel 正在寻找可以升级为主服务器的从服务器。
	"no-good-slave",                     // Sentinel 操作未能找到适合进行升级的从服务器。Sentinel 会在一段时间之后再次尝试寻找合适的从服务器来进行升级，又或者直接放弃执行故障转移操作。
	"selected-slave",                    // Sentinel 顺利找到适合进行升级的从服务器。
	"failover-state-send-slaveof-noone", // Sentinel 正在将指定的从服务器升级为主服务器，等待升级功能完成。
	"failover-end-for-timeout",          // 故障转移因为超时而中止，不过最终所有从服务器都会开始复制新的主服务器（slaves will eventually be configured to replicate with the new master anyway）。
	"failover-end",                      // 故障转移操作顺利完成。所有从服务器都开始复制新的主服务器了。
	"+switch-master",                    // 配置变更，主服务器的 IP 和地址已经改变。 这是绝大多数外部用户都关心的信息。
	"+tilt",                             // 进入 tilt 模式。
	"-tilt",                             // 退出 tilt 模式。
}

// SentinelManagerMonitor 监听哨兵节点的 channel 事件
func SentinelManagerMonitor(addr, password string) {
	errMsg := message.Msg{
		Code:  int64(CodeMonitorSentinelError),
		Title: CodeMonitorSentinelError.Message(),
		Time:  time.Now().Format("2006-01-02 15:04:05"),
	}

	// 创建连接
	rc, err := InitSentinelManagerConn(addr, password)
	if err != nil {
		errMsg.Err = err
		errMsg.Notice()
	}
	defer rc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

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
	for m := range ch {
		//fmt.Println(msg.Channel, msg.Payload, "\r\n")
		// 处理每个 channel 接收到的消息
		switch m.Channel {
		case "+switch-master": // 发生了主从切换
			contentTMP := strings.Split(m.Payload, " ")
			tmp := fmt.Sprintf("master name: %s 发生了主从切换,master 从 %v:%v 切换为 %v:%v",
				contentTMP[0], contentTMP[1], contentTMP[2], contentTMP[3], contentTMP[4])

			errMsg.Err = errors.New(tmp)
			errMsg.Notice()
		default: // 其他哨兵事件通知
			msg := message.Msg{
				Code:    int64(CodeSuccess),
				Title:   "redis 哨兵事件提醒",
				Time:    time.Now().Format("2006-01-02 15:04:05"),
				Content: fmt.Sprintf("%v %v", m.Channel, m.Payload),
			}
			msg.Print()
		}
	}

}

// SentinelStatusCheck 校验哨兵管理节点,master name 对应主从对的状态
//func SentinelStatusCheck(sentinels []string, password, masterName string) {
//	errMsg := message.Msg{
//		Code:  int64(CodeMonitorSentinelError),
//		Title: CodeMonitorSentinelError.Message(),
//		Time:  time.Now().Format("2006-01-02 15:04:05"),
//	}
//
//	// master name 对应 master 实例
//	// 创建连接
//	rm, err := InitSentinelMasterConn(sentinels, password, masterName)
//	if err != nil {
//		errMsg.Err = err
//		errMsg.Notice()
//	}
//	defer rm.Close()
//
//	// master name 对应 slave 实例
//	// 创建连接
//	rs, err := InitSentinelSlaveConn(sentinels, password, masterName)
//	if err != nil {
//		errMsg.Err = err
//		errMsg.Notice()
//	}
//	defer rs.Close()
//
//}
