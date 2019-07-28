package main

import (
	"fmt"
	"github.com/Samanfekri/supernova/queue"
	"time"
)

func main() {
	q := queue.Create("Salam", 1000)
	c1 := queue.Client{Id: "akbar", ReceiveQueue: make(chan interface{}, 100)}
	c2 := queue.Client{Id: "asghar", ReceiveQueue: make(chan interface{}, 100)}
	q.Connect(c1)
	q.Connect(c2)
	go func() {
		for msg := range c1.ReceiveQueue {
			fmt.Println("c1: " + msg.(string))
		}
	}()
	go func() {
		for msg := range c2.ReceiveQueue {
			fmt.Println("c2: " + msg.(string))
		}
	}()
	for i := 0; i < 100; i++ {
		q.Publish(fmt.Sprintf("Publish: %d", i))
		time.Sleep(100 * time.Millisecond)
	}

	q.Broadcast(fmt.Sprintf("Broadcasted message"))
	time.Sleep(100 * time.Millisecond)

	for i := 0; i < 100; i++ {
		err := q.PublishTo(fmt.Sprintf("Publish: %d", i), "asghar")
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
