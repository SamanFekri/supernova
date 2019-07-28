package main

import "errors"

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
