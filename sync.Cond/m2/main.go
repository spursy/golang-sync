package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	var mailbox uint8
	var lock sync.Mutex

	sendCond := sync.NewCond(&lock)
	recvCond := sync.NewCond(&lock)

	send := func(id, index int) {
		lock.Lock()
		for mailbox == 1 {
			sendCond.Wait()
		}
		log.Printf("sender [%d-%d]: the mailbox is empty.", id, index)
		mailbox = 1
		log.Printf("sender [%d-%d]: the letter has been sent.", id, index)
		lock.Unlock()
		recvCond.Broadcast()
	}

	recv := func(id, index int) {
		lock.Lock()
		for mailbox == 0 {
			recvCond.Wait()
		}
		log.Printf("receiver [%d-%d]: the mailbox is full.", id, index)
		mailbox = 0
		log.Printf("receiver [%d-%d]: the letter has been received.", id, index)
		lock.Unlock()
		sendCond.Signal()
	}

	sign := make(chan struct{}, 3)
	max := 6
	go func(id, max int) {
		defer func() {
			sign <- struct{}{}
		}()
		for i := 1; i <= max; i ++ {
			time.Sleep(time.Millisecond * 200)
			send(id, i)
		}
	}(1, max)

	go func(id, max int) {
		defer func() {
			sign <- struct{}{}
		}()
		for j := 1; j <= max; j ++ {
			time.Sleep(time.Millisecond * 200)
			recv(id, j)
		}
	}(2, max/2)

	go func(id, max int) {
		defer func() {
			sign <- struct{}{}
		}()
		for k := 1; k <= max; k ++ {
			time.Sleep(time.Millisecond * 200)
			recv(id, k)
		}
	}(3, max/2)

	<- sign
	<- sign
	<- sign
}