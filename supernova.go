package supernova

import (
	"github.com/Samanfekri/supernova/queue"
)

type Supernova map[string]*Channel

func Dial() *Supernova {
	s := make(Supernova)
	return &s
}

func (s *Supernova) DeclareChannel(name string, size int) *Channel {
	if c, exist := (*s)[name]; exist {
		return c
	}
	(*s)[name] = &Channel{*queue.Create(name, size)}
	return (*s)[name]
}
