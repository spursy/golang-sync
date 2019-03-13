package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	// mailbox 代表信箱
	// 0 代表信箱是空的， 1 代表信箱是满的
	var mailbox uint8 
	// lock 代表信箱上的锁
	var lock sync.RWMutex

	// 代表专用于发信的条件变量
	sendCond := sync.NewCond(&lock)
	// 代表专用于收信的条件变量
	recvCond := sync.NewCond(lock.RLocker())

	// 用于传递演示完成的信号
	sign := make(chan struct{}, 3)
	max := 5

	go func(max int) { // 用于发信
		defer func() {
			sign <- struct{}{}
		}()
		for i := 1; i <= max; i ++ {
			time.Sleep(time.Microsecond * 500)
			lock.Lock()
			for mailbox == 1 {
				sendCond.Wait()
			}
			log.Printf("Sender [%d]: the mailbox is empty.", i)
			mailbox = 1
			log.Printf("Sender [%d]: the letter has been sent.", i)
			lock.Unlock()
			recvCond.Signal()
		}
	}(max)

	go func(max int) { // 用于收信
		defer func() {
			sign <- struct{}{}
		}()

		for j := 1; j <= max; j++ {
			time.Sleep(time.Millisecond * 500)
			lock.RLock()
			for mailbox == 0 {
				recvCond.Wait()
			}
			log.Printf("Receiver [%d]: the mailbox is full.", j)
			mailbox = 0
			log.Printf("Receiver [%d]: the letter has been received.", j)
			lock.RUnlock()
			sendCond.Signal()
		}
	}(max)

	<- sign
	<- sign
}