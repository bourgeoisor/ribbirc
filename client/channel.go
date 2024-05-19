package client

import (
	"fmt"
	"ribbirc/utils"
	"sync"
)

type Channel struct {
	name  string
	topic string

	mutex sync.Mutex
	logs  *utils.Logger
	nicks map[string]bool
}

func newChannel(name string) *Channel {
	return &Channel{
		name:  name,
		logs:  utils.NewLogger(),
		nicks: make(map[string]bool),
	}
}

func (c *Channel) GetLogger() *utils.Logger {
	return c.logs
}

func (c *Channel) UserCount() int {
	return len(c.nicks)
}

func (c *Channel) log(nick string, text string) {
	c.logs.Append(nick, text)
}

func (c *Channel) userMessage(nick string, text string) {
	c.log(nick, text)
}

func (c *Channel) userJoin(nick string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.nicks[nick] = true
	c.log(nick, "JOINED")
}

func (c *Channel) userPart(nick string, reason string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.nicks[nick]; ok {
		delete(c.nicks, nick)
		c.log(nick, fmt.Sprintf("left <%s>", reason))
	}
}

func (c *Channel) userQuit(nick string, reason string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.nicks[nick]; ok {
		delete(c.nicks, nick)
		c.log(nick, fmt.Sprintf("left <%s>", reason))
	}
}
