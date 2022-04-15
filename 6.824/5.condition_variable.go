package main

import "sync"

func main() {
	count, finished := 0, 0
	var mu sync.Mutex
	cond := sync.NewCond(&mu)
	for i := 0; i < 10; i++ {
		go func() {
			vote := requestVote()
			mu.Lock()
			defer mu.Unlock()
			if vote {
				count++
			}
			finished++
			cond.Broadcast() // 如果我们对data做出改变, 就直接call broadcast
		}()
	}

	mu.Lock()
	for count < 5 && finished != 10 {
		cond.Wait() // lock住以后然后等待.
	}
	if count >= 5 {
		println("received 5+ votes")
	} else {
		println("lost")
	}
	mu.Unlock()

}

func requestVote() bool {
	return false
}
