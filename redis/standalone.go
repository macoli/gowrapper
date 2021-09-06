package redisfunc

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// CheckStand 检查单例 redis状态: 是否能连通;是否读写正常
func CheckStand(addr, password string) Msg {
	// 创建 redis 连接
	rc, err := InitStandRedis(addr, password)
	if err != nil {
		return Msg{
			Time:  time.Now().Format("2006-01-02 15:04:05"),
			Code:  CodeMonitorAddrError,
			Title: CodeMonitorAddrError.Message(),
			Err:   err,
		}
	}
	defer rc.Close()

	// 获取 redis 的 info 信息
	infoMap, err := GetInfoMapByClient(rc)
	if err != nil {
		return Msg{
			Time:  time.Now().Format("2006-01-02 15:04:05"),
			Code:  CodeMonitorAddrError,
			Title: CodeMonitorAddrError.Message(),
			Err:   err,
		}
	}

	key := fmt.Sprintf("impossible_exist_key@%v", time.Now().Unix()) // 定义一个 key
	if infoMap["role"] == "master" {                                 // 如果 redis 是 master,则测试写和读
		// 测试写
		err := rc.Set(ctx, key, "test", time.Second*30).Err()
		if err != nil {
			errMsg := fmt.Sprintf("redis 实例(master) %s 测试写数据失败: %v\n", addr, err)
			return Msg{
				Time:  time.Now().Format("2006-01-02 15:04:05"),
				Code:  CodeMonitorAddrError,
				Title: CodeMonitorAddrError.Message(),
				Err:   errors.New(errMsg),
			}
		}

		// 测试读
		err = rc.Get(ctx, key).Err()
		if err != nil && err != redis.Nil {
			errMsg := fmt.Sprintf("redis 实例(master) %s 测试读数据失败: %v\n", addr, err)
			return Msg{
				Time:  time.Now().Format("2006-01-02 15:04:05"),
				Code:  CodeMonitorAddrError,
				Title: CodeMonitorAddrError.Message(),
				Err:   errors.New(errMsg),
			}
		}

	} else if infoMap["role"] == "slave" { // 如果 redis 是 slave,则测试读
		// 测试读
		err = rc.Get(ctx, key).Err()
		if err != nil && err != redis.Nil {
			errMsg := fmt.Sprintf("redis 实例(slave) %s 测试读数据失败: %v\n", addr, err)
			return Msg{
				Time:  time.Now().Format("2006-01-02 15:04:05"),
				Code:  CodeMonitorAddrError,
				Title: CodeMonitorAddrError.Message(),
				Err:   errors.New(errMsg),
			}
		}
	}
	return Msg{
		Code:    CodeSuccess,
		Title:   CodeSuccess.Message(),
		Content: fmt.Sprintf("redis 实例 %s 状态正常", addr),
		Time:    time.Now().Format("2006-01-02 15:04:05"),
	}
}
