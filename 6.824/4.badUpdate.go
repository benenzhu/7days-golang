package main

import (
	"sync"
)

func main() {
	var wg sync.WaitGroup
	counter := 0
	for i := 0; i < 1000; i++ {
		wg.Add(1)
	}
	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done() // not ensure the crtical section to right
			counter++
		}()
	}
	wg.Wait()
	println(counter) // this will not be 1000
}
