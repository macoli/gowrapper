package redisfunc

import (
	"sync"
	"time"
)

var (
	redisWait        *sync.WaitGroup
	sentinelChannels = []string{
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
)

type ClusterNode struct {
	ID          string   // 当前节点 ID
	Addr        string   // 当前节点地址(ip:port)
	ClusterPort string   // 当前节点和集群其他节点通信端口(默认为节点端口+1000),3.x 版本不展示该信息
	Flags       []string // 当前节点标志:myself, master, slave, fail?, fail, handshake, noaddr, nofailover, noflags
	MasterID    string   // 如果当前节点是 slave,这里就是 对应 master 的 ID,如果当前节点是 master,以"-"表示
	PingSent    int64    // 最近一次发送ping的时间，这个时间是一个unix毫秒时间戳，0代表没有发送过
	PongRecv    int64    // 最近一次收到pong的时间，使用unix时间戳表示
	ConfigEpoch int64    // 节点的epoch值.每当节点发生失败切换时，都会创建一个新的，独特的，递增的epoch。如果多个节点竞争同一个哈希槽时，epoch值更高的节点会抢夺到。
	LinkState   string   // node-to-node集群总线使用的链接的状态: connected或disconnected
	Slots       []string // 哈希槽值或者一个哈希槽范围
}

type MasterSlaveMap struct {
	MasterID   string
	MasterAddr string
	SlaveAddr  string
	SlaveID    string
	SlotStr    string
}

type ClusterInfo struct {
	ClusterNodes    []*ClusterNode
	MasterSlaveMaps []*MasterSlaveMap
	Masters         []string
	Slaves          []string
	IDToAddr        map[string]string
	AddrToID        map[string]string
}

type SlowLog struct {
	Instance string
	Command  string
	Duration time.Duration
	Time     string
}

type AppStandAlone struct {
	AppID     int64
	AppName   string
	Password  string
	Masters   []string
	Slaves    []string
	Instances []string
}

type AppSentinel struct {
	AppID       int64
	AppName     string
	Password    string
	Sentinels   []string
	MasterNames []string
}

type AppCluster struct {
	AppID     int64
	AppName   string
	Password  string
	Instances []string
}

type App struct {
	Type    string
	AppInfo interface{}
}
