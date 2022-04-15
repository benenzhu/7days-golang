package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan bool, 1)
	go func() {
		time.Sleep(200 * time.Millisecond)
		<-c
	}()
	start := time.Now()
	c <- true
	fmt.Printf("send took %v\n", time.Since(start))
}
