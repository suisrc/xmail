package core

import (
	"bufio"
	"os"
	"sync"
)

// 线程并发执行任务
func RunByTaskByIdx(task func(idx int), count, cmax int) {
	wg := sync.WaitGroup{}
	wk := make(chan int, cmax)
	for i := 0; i < count; i++ {
		wk <- i // 限制并发
		wg.Add(1)
		go func(idx int) {
			defer func() {
				<-wk // 释放
				wg.Done()
			}()
			task(idx)
		}(i)
	}
	wg.Wait()
}

// 追加写入文件, os,WireteFile
func WriteFileAppend(path string, data []byte) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	return err
}

// 读取文件, os.ReadFile
func ReadFileLines(path string) ([]string, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
