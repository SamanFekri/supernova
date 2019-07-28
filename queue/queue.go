package queue

import (
	"errors"
	"math/rand"
	"time"
)

type Queue struct {
	Name        string
	MaxCapacity int
	Q           chan message
	Clients     map[string]Client
	w           chan bool
}

type message struct {
	isBroadcast bool
	receiver    string
	body        interface{}
}

func Create(Name string, MaxCapacity int) *Queue {
	q := &Queue{
		Name:        Name,
		MaxCapacity: MaxCapacity,
		Q:           make(chan message, MaxCapacity),
		Clients:     make(map[string]Client),
		w:           make(chan bool),
	}
	go q.listen()
	return q
}

func (q *Queue) Connect(client Client) {
	q.Clients[client.Id] = client
	if len(q.Clients) > 0 {
		q.w <- true
	}
}

func (q *Queue) Disconnect(id string) {
	delete(q.Clients, id)
}

func (q *Queue) Publish(input interface{}) {
	q.Q <- message{isBroadcast: false, receiver: "", body: input}
}

func (q *Queue) PublishTo(input interface{}, receiver string) error {
	if _, exist := q.Clients[receiver]; !exist {
		return errors.New("This receiver does`nt exist in this queue.")
	}
	q.Q <- message{isBroadcast: false, receiver: receiver, body: input}
	return nil
}

func (q *Queue) Broadcast(input interface{}) {
	q.Q <- message{isBroadcast: true, receiver: "", body: input}
}

func (q *Queue) listen() {
	for m := range q.Q {
		if len(q.Clients) < 1 {
			<-q.w
		}
		if m.isBroadcast { // send the message to all receivers
			for _, client := range q.Clients {
				go client.publishToClient(m.body)
			}
		} else if m.receiver != "" { // send the message to specific receiver
			if client, exist := q.Clients[m.receiver]; exist {
				go client.publishToClient(m.body)
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
			go client.publishToClient(m.body)
		}
	}
}
