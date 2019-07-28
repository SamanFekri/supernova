package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Queue struct {
	Name        string
	MaxCapacity int
	Queue       chan message
	Clients     map[string]Client
}

type message struct {
	isBroadcast bool
	receiver    string
	body        interface{}
}

type Client struct {
	Id           string
	ReceiveQueue chan interface{}
}

func Create(Name string, MaxCapacity int) *Queue {
	return &Queue{
		Name:        Name,
		MaxCapacity: MaxCapacity,
		Queue:       make(chan message, MaxCapacity),
		Clients:     make(map[string]Client),
	}
}

func (q *Queue) Publish(input interface{}) {
	q.Queue <- message{isBroadcast: false, receiver: "", body: input}
}

func (q *Queue) PublishTo(input interface{}, receiver string) error {
	if _, exist := q.Clients[receiver]; !exist {
		return errors.New("This receiver does`nt exist in this queue.")
	}
	q.Queue <- message{isBroadcast: false, receiver: receiver, body: input}
	return nil
}

func (q *Queue) Broadcast(input interface{}) {
	q.Queue <- message{isBroadcast: true, receiver: "", body: input}
}

func (q *Queue) listen() {
	for m := range q.Queue {
		if m.isBroadcast { // send the message to all receivers
			for _, client := range q.Clients {
				go publishToClient(m.body, client.ReceiveQueue)
			}
		} else if m.receiver != "" { // send the message to specific receiver
			if client, exist := q.Clients[m.receiver]; exist {
				go publishToClient(m.body, client.ReceiveQueue)
			}
		} else { // send the message to a random receiver
			n := len(q.Clients)
			rand.Seed(time.Now().UnixNano())
			i := rand.Intn(n)
			var client Client
			for _, client = range q.Clients {
				if i == 0 {
					break
				}
				i--
			}
			go publishToClient(m.body, client.ReceiveQueue)
		}
	}
}

func publishToClient(input interface{}, q chan interface{}) {
	fmt.Println(input)
	q <- input
}
