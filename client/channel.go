package client

import (
	"fmt"
	"ribbirc/utils"
	"sync"
)

type Channel struct {
	Name  string
	Topic string

	mutex sync.Mutex
	Logs  *utils.Logger
	Nicks map[string]bool
}

func newChannel(name string) *Channel {
	return &Channel{
		Name:  name,
		Logs:  utils.NewLogger(),
		Nicks: make(map[string]bool),
	}
}

func (c *Channel) log(nick string, text string) {
	c.Logs.Append(nick, text)
}

func (c *Channel) userMessage(nick string, text string) {
	c.log(nick, text)
}

func (c *Channel) userJoin(nick string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.Nicks[nick] = true
	// @todo: handle prefixes
	c.log(nick, "JOINED")
}

func (c *Channel) usersJoin(nicks []string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, nick := range nicks {
		c.Nicks[nick] = true
		// @todo: handle prefixes
	}
}

func (c *Channel) userPart(nick string, reason string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.Nicks[nick]; ok {
		delete(c.Nicks, nick)
		c.log(nick, fmt.Sprintf("left <%s>", reason))
	}
}

func (c *Channel) userQuit(nick string, reason string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.Nicks[nick]; ok {
		delete(c.Nicks, nick)
		c.log(nick, fmt.Sprintf("left <%s>", reason))
	}
}
