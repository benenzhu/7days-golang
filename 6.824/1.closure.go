package main

import "sync"

func main() {
	var a string
	var wg sync.WaitGroup // 创建
	wg.Add(1)             //  + 1
	go func() {
		a = "hello world"
		wg.Done() // - 1
	}()
	wg.Wait() // 这里需要等待
	println(a)
}
