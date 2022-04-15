package main

import "sync"

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			sendRpc(i) // 这样永远都是 5 会提醒的吧...
			wg.Done()  // loopclosure
		}()
	}
	wg.Wait()
}
func sendRpc(x int) {
	println(x)
}
