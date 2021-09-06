package redisfunc

import (
	"math/rand"
	"sync"
	"time"
)

var monitorWG *sync.WaitGroup

type M interface {
	Type()
	Run(wg *sync.WaitGroup)
}

type MonitorQueueChannel struct {
	RedisTimer chan App
	RedisMsg   chan Msg
}

func (c MonitorQueueChannel) Type() string {
	return "channel"
}

// NewMonitorQueueWithChannel 初始化消息对列(使用 channel)
func NewMonitorQueueWithChannel(timerSize, msgSize int64) *MonitorQueueChannel {
	return &MonitorQueueChannel{
		RedisTimer: make(chan App, timerSize),
		RedisMsg:   make(chan Msg, msgSize),
	}
}

// Run 启动函数
func (c *MonitorQueueChannel) Run(wg *sync.WaitGroup) {
	monitorWG = wg
	monitorWG.Add(3)
	go GetRedisAppInfoTimerWithChannel(c.RedisTimer)
	go CheckRedisAppInfoWithChannel(c.RedisTimer, c.RedisMsg)
	go DealRespMsgWithChannel(c.RedisMsg)
}

// GetRedisAppInfoTimerWithChannel 使用定时器定时向 redisTimer channel 中发送要监控的 redis 应用
func GetRedisAppInfoTimerWithChannel(redisTimer chan<- App) {
	defer monitorWG.Done()
	// 设置定时器
	rand.Seed(time.Now().UnixNano())
	interval := time.Second * time.Duration(0+rand.Intn(6)+1)
	timer := time.NewTicker(interval) // 随机 60-120 s

	// 定时执行任务,获取结果并处理
	for _ = range timer.C {
		// 获取 app 信息并发送到 RedisTimer channel 中
	}

}

// CheckRedisAppInfoWithChannel 从 redisTimer channel 中获取 redis 应用, 并检测其状态,将结果发送到 redisMsg channel 中
func CheckRedisAppInfoWithChannel(redisTimer <-chan App, redisMsg chan<- Msg) {
	defer monitorWG.Done()
	// 并发 redisTimer 数量开始从 redisTimer channel 中获取信息并处理
	//for i = 0; i < cap(redisTimer); i++ {
	//
	//}
}

// DealRespMsgWithChannel 从 redisMsg channel 中获取检测结果,并处理(打印/发送第三方)
func DealRespMsgWithChannel(redisMsg <-chan Msg) {
	defer monitorWG.Done()
}

// MonitorStandalone 监控 redis 实例
// 使用定时器,每个定时器时间间隔:60-120s 随机
// 监控功能: redis 是否连接正常(建立连接+ping); master 读写功能; slave 读功能
func MonitorStandalone(wg *sync.WaitGroup, app AppStandAlone) {
	defer wg.Done()

	// 设置定时器
	rand.Seed(time.Now().UnixNano())
	interval := time.Second * time.Duration(0+rand.Intn(6)+1)
	timer := time.NewTicker(interval) // 随机 60-120 s

	// 定时执行任务,获取结果并处理
	for _ = range timer.C {
		for _, addr := range app.Instances {
			msg := CheckStand(addr, app.Password)
			msg.Notice()
		}

	}
}

// MonitorSentinel 监控 redis 哨兵
// 监控功能: 哨兵节点是否连接正常(建立连接+ping); 哨兵事件; 哨兵对应主从对监控(同单实例 redis 监控)
func MonitorSentinel(wg *sync.WaitGroup, addr, password string, app AppSentinel) {
	defer wg.Done()

	// 定时连接所有哨兵节点,判断节点是否连接正常

	// 哨兵管理节点监控
	for _, sentinel := range app.Sentinels {
		CheckSentinelManager(sentinel, app.Password)
	}

	// 哨兵主从对监控

}

// MonitorCluster 监控 redis 集群
// 监控功能: 集群节点是否连接正常(建立连接+ping); 集群状态是否正常; 集群事件
func MonitorCluster(wg *sync.WaitGroup, app AppCluster) {

}

// =====================================================================================
type MonitorQueueRedis struct {
	Addr     string
	Password string
}

func (r MonitorQueueRedis) Type() string {
	return "redis"
}

// NewMonitorQueueWithRedis 初始化消息对列(使用 redis)
func NewMonitorQueueWithRedis(addr, password string) *MonitorQueueRedis {
	return &MonitorQueueRedis{
		Addr:     addr,
		Password: password,
	}
}
