package main

import (
	"fmt"
	"github.com/Samanfekri/supernova"
	"github.com/Samanfekri/supernova/queue"
	"time"
)

func main() {
	s := supernova.Dial()
	ch := s.DeclareChannel("channel1", 1000)

	c1 := queue.Client{Id: "abc", ReceiveQueue: make(chan interface{}, 100)}
	c2 := queue.Client{Id: "cba", ReceiveQueue: make(chan interface{}, 100)}

	u1 := ch.Connect(c1)
	u2 := ch.Connect(c2)

	for i := 0; i < 100; i++ {
		u1.Publish(fmt.Sprintf("C1 send to channel1 %d", i))

	}

	go func() {
		for msg := range u2.Receive() {
			fmt.Println("Receive: ", msg)
		}
	}()

	time.Sleep(50 * time.Millisecond)
}
