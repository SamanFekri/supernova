package supernova

import "github.com/Samanfekri/supernova/queue"

type Channel struct {
	queue.Queue
}

func (c *Channel) Connect(u queue.Client) *ConnectedUser {
	return &ConnectedUser{User: u, Channel: *c}
}

type ConnectedUser struct {
	User    queue.Client
	Channel Channel
}

func (u *ConnectedUser) Publish(msg interface{}) {
	u.Channel.Queue.Publish(msg)
}

func (u *ConnectedUser) PublishTo(msg interface{}, receiver string) error {
	err := u.Channel.Queue.PublishTo(msg, receiver)
	return err
}

func (u *ConnectedUser) Broadcast(msg interface{}) {
	u.Channel.Queue.Broadcast(msg)
}

func (u *ConnectedUser) Receive() chan interface{} {
	return u.User.ReceiveQueue
}
