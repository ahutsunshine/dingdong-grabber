package main

import (
	"sync"
	"time"
)

var (
	waitGroup sync.WaitGroup
)

func main() {

	go ddMain(&waitGroup)
	go mtMain(&waitGroup)

	time.Sleep(3 * time.Second)
	waitGroup.Wait()

}
