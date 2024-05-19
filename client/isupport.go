package client

import (
	"strconv"
	"strings"
)

type ISupport struct {
	awaylen     int
	casemapping string
	chanlimit   string
	chanmodes   string
	channellen  int
	chantypes   string
	elist       string
	excepts     string
	extban      string
	hostlen     int
	invex       string
	kicklen     int
	maxlist     string
	maxtargets  int
	modes       int
	network     string
	nicklen     int
	prefix      string
	safelist    bool
	silence     int
	statusmsg   string
	targmax     string
	topiclen    int
	userlen     int
}

func newISupport() *ISupport {
	return &ISupport{
		chantypes: "#&",
		excepts:   "e",
		invex:     "I",
	}
}

func (i *ISupport) parseRpl(tokens []string) {
	for _, token := range tokens {
		i.parseToken(token)
	}
}

func (i *ISupport) parseToken(token string) {
	name := token
	value := ""
	equalSign := strings.Index(token, "=")
	if equalSign > 0 {
		name = token[:equalSign]
		value = token[equalSign+1:]
	}

	switch name {
	case "-AWAYLEN":
		i.awaylen = 0
	case "AWAYLEN":
		i.awaylen, _ = strconv.Atoi(value)
	case "-CASEMAPPING":
		i.casemapping = ""
	case "CASEMAPPING":
		i.casemapping = value
	case "-CHANLIMIT":
		i.chanlimit = ""
	case "CHANLIMIT":
		i.chanlimit = value
	case "-CHANMODES":
		i.chanmodes = ""
	case "CHANMODES":
		i.chanmodes = value
	case "-CHANNELLEN":
		i.channellen = 0
	case "CHANNELLEN":
		i.channellen, _ = strconv.Atoi(value)
	case "-CHANTYPES":
		i.chantypes = "#&"
	case "CHANTYPES":
		i.chantypes = value
	case "-ELIST":
		i.elist = ""
	case "ELIST":
		i.elist = value
	case "-EXCEPTS":
		i.excepts = "e"
	case "EXCEPTS":
		i.excepts = value
	case "-EXTBAN":
		i.extban = ""
	case "EXTBAN":
		i.extban = value
	case "-HOSTLEN":
		i.hostlen = 0
	case "HOSTLEN":
		i.hostlen, _ = strconv.Atoi(value)
	case "-INVEX":
		i.invex = "I"
	case "INVEX":
		i.invex = value
	case "-KICKLEN":
		i.kicklen = 0
	case "KICKLEN":
		i.kicklen, _ = strconv.Atoi(value)
	case "-MAXLIST":
		i.maxlist = ""
	case "MAXLIST":
		i.maxlist = value
	case "-MAXTARGETS":
		i.maxtargets = 0
	case "MAXTARGETS":
		i.maxtargets, _ = strconv.Atoi(value)
	case "-MODES":
		i.modes = 0
	case "MODES":
		i.modes, _ = strconv.Atoi(value)
	case "-NETWORK":
		i.network = ""
	case "NETWORK":
		i.network = value
	case "-NICKLEN":
		i.nicklen = 0
	case "NICKLEN":
		i.nicklen, _ = strconv.Atoi(value)
	case "-PREFIX":
		i.prefix = ""
	case "PREFIX":
		i.prefix = value
	case "-SAFELIST":
		i.safelist = false
	case "SAFELIST":
		i.safelist = true
	case "-SILENCE":
		i.silence = 0
	case "SILENCE":
		i.silence, _ = strconv.Atoi(value)
	case "-STATUSMSG":
		i.statusmsg = ""
	case "STATUSMSG":
		i.statusmsg = value
	case "-TARGMAX":
		i.targmax = ""
	case "TARGMAX":
		i.targmax = value
	case "-TOPICLEN":
		i.topiclen = 0
	case "TOPICLEN":
		i.topiclen, _ = strconv.Atoi(value)
	case "-USERLEN":
		i.userlen = 0
	case "USERLEN":
		i.userlen, _ = strconv.Atoi(value)
	}
}
