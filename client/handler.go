package client

import (
	"fmt"
	"ribbirc/utils"
	"strconv"
	"strings"
	"time"
)

func (s *Server) handleServerMessage(message *utils.Message) {
	switch message.Command {
	case "NOTICE":
		s.log(message.Parameters[1])

	case "PING":
		s.SendMessage(&utils.Message{Command: "PONG"})

	case "PONG":
		then, _ := strconv.ParseInt(message.Parameters[1], 10, 64)
		diff := time.Now().UnixMilli() - then
		s.log(fmt.Sprintf("Pong received, response time %dms", diff))

	case "ERROR":
		s.log(message.Parameters[0])

	case "JOIN":
		if message.SourceNick() == s.nick {
			s.channelsJoined[message.Parameters[0]] = newChannel(message.Parameters[0])
			s.channelsJoined[message.Parameters[0]].userJoin(s.nick)
		} else {
			s.channelsJoined[message.Parameters[0]].userJoin(message.SourceNick())
		}

	case "PART":
		if message.SourceNick() == s.nick {
			delete(s.channelsJoined, message.Parameters[0])
		} else {
			reason := ""
			if len(message.Parameters) > 1 {
				reason = message.Parameters[1]
			}
			s.channelsJoined[message.Parameters[0]].userPart(message.SourceNick(), reason)
		}

	case "QUIT":
		if message.SourceNick() == s.nick {
			// @todo: exit server
		} else {
			reason := ""
			if len(message.Parameters) > 0 {
				reason = message.Parameters[0]
			}
			for channel, _ := range s.channelsJoined {
				s.channelsJoined[channel].userQuit(message.SourceNick(), reason)
			}
		}

	case "NICK":
		if message.SourceNick() == s.nick {
			s.nick = message.Parameters[0]
		}
		for channel, _ := range s.channelsJoined {
			s.channelsJoined[channel].userNick(message.SourceNick(), message.Parameters[0])
		}

	case "PRIVMSG":
		s.channelsJoined[message.Parameters[0]].userMessage(message.SourceNick(), message.Parameters[1])

	case utils.RPL_WELCOME:
		// <client> :Welcome to the <networkname> Network, <nick>[!<user>@<host>]
		s.log(message.Parameters[1])

	case utils.RPL_YOURHOST:
		// <client> :Your host is <servername>, running version <version>
		s.log(message.Parameters[1])

	case utils.RPL_CREATED:
		// <client> :This server was created <datetime>
		s.log(message.Parameters[1])

	case utils.RPL_MYINFO:
		// <client> <servername> <version> <available user modes>
		// <available channel modes> [<channel modes with a parameter>]
		s.name = message.Parameters[1]
		s.version = message.Parameters[2]
		s.availableServerModes = message.Parameters[3]
		s.availableChannelModes = message.Parameters[4]

	case utils.RPL_ISUPPORT:
		// <client> <1-13 tokens> :are supported by this server
		s.iSupport.parseRpl(message.Parameters[1 : len(message.Parameters)-1])

	case utils.RPL_LUSERCLIENT:
		// <client> :There are <u> users and <i> invisible on <s> servers
		s.log(message.Parameters[1])

	case utils.RPL_LUSEROP:
		// <client> <ops> :operator(s) online
		s.log(fmt.Sprintf("%s %s", message.Parameters[1], message.Parameters[2]))

	case utils.RPL_LUSERUNKNOWN:
		// <client> <connections> :unknown connection(s)
		s.log(fmt.Sprintf("%s %s", message.Parameters[1], message.Parameters[2]))

	case utils.RPL_LUSERCHANNELS:
		// <client> <channels> :channels formed
		s.log(fmt.Sprintf("%s %s", message.Parameters[1], message.Parameters[2]))

	case utils.RPL_LUSERME:
		// <client> :I have <c> clients and <s> servers
		s.log(message.Parameters[1])

	case utils.RPL_LOCALUSERS:
		// <client> [<u> <m>] :Current local users <u>, max <m>
		s.log(message.Parameters[len(message.Parameters)-1])

	case utils.RPL_GLOBALUSERS:
		// <client> [<u> <m>] :Current global users <u>, max <m>
		s.log(message.Parameters[len(message.Parameters)-1])

	case utils.RPL_NOTOPIC:
		// <client> <channel> :No topic is set
		s.channelsJoined[message.Parameters[1]].Topic = ""

	case utils.RPL_TOPIC:
		// <client> <channel> :<topic>
		s.channelsJoined[message.Parameters[1]].Topic = message.Parameters[2]
		s.channelsJoined[message.Parameters[1]].Logs.Append("*", utils.LogSystem, message.Parameters[2])

	case utils.RPL_TOPICWHOTIME:
		//<client> <channel> <nick> <setat>
		seconds, _ := strconv.ParseInt(message.Parameters[3], 10, 64)
		text := fmt.Sprintf("Topic set by %s on %s.", message.ParamNick(2), time.Unix(seconds, 0))
		s.channelsJoined[message.Parameters[1]].Logs.Append("*", utils.LogSystem, text)

	case utils.RPL_NAMREPLY:
		// <client> <symbol> <channel> :[prefix]<nick>{ [prefix]<nick>}
		s.channelsJoined[message.Parameters[2]].usersJoin(strings.Split(message.Parameters[3], " "))

	case utils.RPL_ENDOFNAMES:
		// <client> <channel> :End of /NAMES list
		break

	case utils.RPL_MOTDSTART:
		// <client> :- <server> Message of the day -
		s.motd = append(s.motd, message.Parameters[1])

	case utils.RPL_MOTD:
		// <client> :<line of the motd>
		s.motd = append(s.motd, message.Parameters[1])

	case utils.RPL_ENDOFMOTD:
		// <client> :End of /MOTD command.
		for _, m := range s.motd {
			s.log(m)
		}

	default:
		text := fmt.Sprintf("Unimplemented reply: %s", utils.MarshalMessage(message))
		s.logs.Append("System", utils.LogError, text)
	}
}
