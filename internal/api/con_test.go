package api

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

func Test_Con(t *testing.T) {
	results := []string{"1", "2", "3", "4", "5", "6", "7"}
	// 定义 goroutines 等待组
	var wg sync.WaitGroup
	// 并发数量
	maxConcurrency := 3
	// 创建 sem 通道
	sem := make(chan struct{}, maxConcurrency)
	// 用于存储文件地址的通道
	filePaths := make(chan string, len(results))
	for i, item := range results {
		log.Printf("当前下标：[%d]\n", i+1)
		log.Printf("当前资源：[%s]\n", item)
		wg.Add(1)
		go down(item, sem, &wg, filePaths)
	}
	// 等待所有下载完成
	go func() {
		wg.Wait()
		close(filePaths)
	}()
	for f := range filePaths {
		fmt.Println("处理file:", f)
	}
}

func down(urlStr string, sem chan struct{}, wg *sync.WaitGroup, filePaths chan string) {
	defer wg.Done()

	// 从信号量中获取一个令牌
	sem <- struct{}{}
	defer func() { <-sem }() // 确保在函数返回时释放信号量令牌

	fmt.Println("下载中", urlStr)
	time.Sleep(1 * time.Second)
	fmt.Println("下载完成", urlStr)

	// 发送文件路径到通道
	filePaths <- urlStr

}
