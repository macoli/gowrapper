package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

//InitStandConn 初始化单例 redis 连接
func InitStandConn(addr, password string) (*redis.Client, error) {
	rc := redis.NewClient(&redis.Options{
		Addr:        addr,
		Password:    password,
		DB:          0,
		PoolSize:    100,
		DialTimeout: time.Minute * 30,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rc.Ping(ctx).Result()
	if err != nil {
		errMsg := fmt.Sprintf("redis 实例 %s 连接失败: %v\n", addr, err)
		return nil, errors.New(errMsg)
	}
	return rc, nil
}

//InitSentinelMasterConn 初始化哨兵连接,通过哨兵获取到对应 master name 节点的 master 连接
func InitSentinelMasterConn(addrSlice []string, password, masterName string) (*redis.Client, error) {
	rc := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    masterName,
		SentinelAddrs: addrSlice,
		Password:      password,
		PoolSize:      1000,
		DialTimeout:   time.Minute * 30,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := rc.Ping(ctx).Result()
	if err != nil {
		errMsg := fmt.Sprintf("哨兵 %v 上 %s 的 master 连接失败: %v\n", addrSlice, masterName, err)
		return nil, errors.New(errMsg)
	}
	return rc, nil
}

//InitSentinelSlaveConn 初始化哨兵连接,通过哨兵获取到对应 master name 节点的 slave 只读连接
func InitSentinelSlaveConn(addrSlice []string, password, masterName string) (*redis.ClusterClient, error) {
	rc := redis.NewFailoverClusterClient(&redis.FailoverOptions{
		MasterName:    masterName,
		SentinelAddrs: addrSlice,
		Password:      password,
		PoolSize:      1000,
		DialTimeout:   time.Minute * 30,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rc.Ping(ctx).Result()
	if err != nil {
		errMsg := fmt.Sprintf("哨兵 %v 上 %s 的 slave 连接失败: %v\n", addrSlice, masterName, err)
		return nil, errors.New(errMsg)
	}
	return rc, nil
}

//InitSentinelManagerConn 初始化哨兵管理连接,用于连接哨兵节点,管理哨兵
func InitSentinelManagerConn(addr, password string) (*redis.SentinelClient, error) {
	rc := redis.NewSentinelClient(&redis.Options{
		Addr:     addr,
		Password: password,
		PoolSize: 100,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rc.Ping(ctx).Result()
	if err != nil {
		errMsg := fmt.Sprintf("哨兵管理节点: %s 连接失败: %v\n", addr, err)
		return nil, errors.New(errMsg)
	}
	return rc, nil
}

//InitClusterConn 初始化集群连接
func InitClusterConn(addrSlice []string, password string) (*redis.ClusterClient, error) {
	rc := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    addrSlice,
		Password: password,
		PoolSize: 1000,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rc.Ping(ctx).Result()
	if err != nil {
		errMsg := fmt.Sprintf("集群节点: %s 连接失败: %v\n", addrSlice, err)
		return nil, errors.New(errMsg)
	}
	return rc, nil
}
