package queue

type Client struct {
	Id           string
	ReceiveQueue chan interface{}
}

func (c *Client) publishToClient(input interface{}) {
	c.ReceiveQueue <- input
}
