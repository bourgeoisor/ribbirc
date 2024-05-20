package client

import (
	"fmt"
	"ribbirc/utils"
	"strings"
)

func (s *Server) handleServerMessage(message *utils.Message) {
	switch message.Command {
	case "NOTICE":
		s.log(message.Parameters[1])
	case "PING":
		s.SendMessage(&utils.Message{Command: "PONG"})
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
	case "PRIVMSG":
		s.channelsJoined[message.Parameters[0]].userMessage(message.SourceNick(), message.Parameters[1])
	case utils.RPL_WELCOME:
		s.log(message.Parameters[1])
	case utils.RPL_YOURHOST:
		s.log(message.Parameters[1])
	case utils.RPL_CREATED:
		s.log(message.Parameters[1])
	case utils.RPL_MYINFO:
		s.name = message.Parameters[1]
		s.version = message.Parameters[2]
		s.availableServerModes = message.Parameters[3]
		s.availableChannelModes = message.Parameters[4]
	case utils.RPL_ISUPPORT:
		s.iSupport.parseRpl(message.Parameters[1 : len(message.Parameters)-1])
	case utils.RPL_LUSERCLIENT:
		s.log(message.Parameters[1])
	case utils.RPL_LUSEROP:
		s.log(fmt.Sprintf("%s %s", message.Parameters[1], message.Parameters[2]))
	case utils.RPL_LUSERUNKNOWN:
		s.log(fmt.Sprintf("%s %s", message.Parameters[1], message.Parameters[2]))
	case utils.RPL_LUSERCHANNELS:
		s.log(fmt.Sprintf("%s %s", message.Parameters[1], message.Parameters[2]))
	case utils.RPL_LUSERME:
		s.log(message.Parameters[1])
	case utils.RPL_LOCALUSERS:
		s.log(message.Parameters[len(message.Parameters)-1])
	case utils.RPL_GLOBALUSERS:
		s.log(message.Parameters[len(message.Parameters)-1])
	case utils.RPL_NOTOPIC: // <client> <channel> :No topic is set
		s.channelsJoined[message.Parameters[1]].Topic = ""
	case utils.RPL_TOPIC: // <client> <channel> :<topic>
		s.channelsJoined[message.Parameters[1]].Topic = message.Parameters[2]
		s.channelsJoined[message.Parameters[1]].Logs.Append("*", utils.LogSystem, message.Parameters[2])
	case utils.RPL_NAMREPLY: // <client> <symbol> <channel> :[prefix]<nick>{ [prefix]<nick>}
		s.channelsJoined[message.Parameters[2]].usersJoin(strings.Split(message.Parameters[3], " "))
	case utils.RPL_ENDOFNAMES: // <client> <channel> :End of /NAMES list
		break
	case utils.RPL_MOTDSTART:
		s.motd = append(s.motd, message.Parameters[1])
	case utils.RPL_MOTD:
		s.motd = append(s.motd, message.Parameters[1])
	case utils.RPL_ENDOFMOTD:
		for _, m := range s.motd {
			s.log(m)
		}
	default:
		s.log(fmt.Sprintf("Unimplemented RPL: %s", utils.MarshalMessage(message)))
	}
}
