package sample

import (
	"context"
	"fmt"
	"sync"

	"github.com/macoli/iwrapper/iredis"
	"github.com/macoli/iwrapper/istring"
)

// ClusterWriteTestData 并发向集群中写测试数据
// 并发数默认为 10
// 写入 1000w key 名为: key%i ; value长度为 128 字节
func ClusterWriteTestData(addrList []string, password string, keyNums, workerNums int) {
	var wg sync.WaitGroup
	var workerChannel chan struct{}
	if workerNums > 0 {
		workerChannel = make(chan struct{}, workerNums)
	} else {
		workerChannel = make(chan struct{}, 10)
	}

	rc, err := iredis.InitClusterConn(addrList, password)
	if err != nil {
		fmt.Println(err)
	}
	defer rc.Close()

	// 向 dataChannel 生产数据
	for i := 0; i < keyNums; i++ {
		workerChannel <- struct{}{} // 添加信号,当 workerChannel 满了之后就会阻塞创建新的 goroutine
		wg.Add(1)

		go func(x int) {
			defer wg.Done()
			// 执行完毕,释放信号
			defer func() {
				<-workerChannel
			}()

			key := fmt.Sprintf("key%d", x)
			value := istring.RandString(128)
			rc.Set(context.Background(), key, value, 0)
		}(i)
	}
	wg.Wait()
}
