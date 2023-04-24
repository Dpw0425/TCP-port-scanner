package main

import (
	"flag"
	"fmt"
	"sync"
	"tcp-port-scanner/pkg/logger"
	"tcp-port-scanner/pkg/utils"
	"time"
)

func main() {
	// 命令行参数解析
	hostName := flag.String("hostName", "127.0.0.1", "hostname to test")
	startPort := flag.Int("startPort", 1, "the port on which the scanning starts")
	endPort := flag.Int("endPort", 8080, "the port from which the scanning ends")
	timeOut := flag.Duration("timeOut", time.Millisecond*200, "timeout")
	maxConcurrency := flag.Int("maxConcurrency", 1000, "maximum supported concurrency")
	flag.Parse()

	logger.InitZap() // 初始化日志

	ports := []int{} // 运行结果保存到此数组

	wg := &sync.WaitGroup{}
	mutex := &sync.Mutex{}                    // 设置互斥锁，保证同一时间有且只有一个 goroutine 进入临界区
	c := make(chan struct{}, *maxConcurrency) // 限制同时开启的最大线程数
	defer close(c)

	for port := *startPort; port <= *endPort; port++ {
		wg.Add(1) // 开启新的线程，计数器 +1
		c <- struct{}{}
		// 开启 goroutine
		go func(port int) {
			status := utils.Scanner(*hostName, port, *timeOut) // 设置等待时长，避免程序运行过快接收不到返回的错误信息
			if status {
				mutex.Lock()
				ports = append(ports, port)
				mutex.Unlock()
			}
			defer wg.Done() // 任务完成，计数器 -1
			<-c
		}(port)
	}

	wg.Wait() // 等待所有并发任务执行完成

	fmt.Printf("opened ports: %v\n", ports)
	logger.Info("opened ports: ", ports) // 写入日志
}
