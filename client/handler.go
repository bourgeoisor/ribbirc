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

	case utils.RPL_STATSCOMMANDS:
		// <client> <command> <count> [<byte count> <remote count>]
		text := strings.Join(message.Parameters[1:], " ")
		s.BufferStats = append(s.BufferStats, text)

	case utils.RPL_ENDOFSTATS:
		// <client> <stats letter> :End of /STATS report
		for _, m := range s.BufferStats {
			s.log(fmt.Sprintf("[stats] %s", m))
		}
		s.BufferStats = make([]string, 0)

	case utils.RPL_STATSUPTIME:
		// <client> :Server Up <days> days <hours>:<minutes>:<seconds>
		text := strings.Join(message.Parameters[1:], " ")
		s.BufferStats = append(s.BufferStats, text)

	case utils.RPL_STATSCONN:
		// :Highest connection count: %d (%d clients) (%lu connections received)
		s.log(message.Parameters[1])

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

	case utils.RPL_ADMINME:
		// <client> [<server>] :Administrative info
		s.log(message.Parameters[len(message.Parameters)-1])

	case utils.RPL_ADMINLOC1:
		// <client> :<info>
		s.log(message.Parameters[1])

	case utils.RPL_ADMINLOC2:
		// <client> :<info>
		s.log(message.Parameters[1])

	case utils.RPL_ADMINEMAIL:
		// <client> :<info>
		s.log(message.Parameters[1])

	case utils.RPL_TRYAGAIN:
		// <client> <command> :Please wait a while and try again.
		s.log(message.Parameters[2])

	case utils.RPL_LOCALUSERS:
		// <client> [<u> <m>] :Current local users <u>, max <m>
		s.log(message.Parameters[len(message.Parameters)-1])

	case utils.RPL_GLOBALUSERS:
		// <client> [<u> <m>] :Current global users <u>, max <m>
		s.log(message.Parameters[len(message.Parameters)-1])

	case utils.RPL_WHOISCERTFP:
		// <client> <nick> :has client certificate fingerprint <fingerprint>
		text := strings.Join(message.Parameters[1:], " ")
		s.BufferWho = append(s.BufferWho, text)

	case utils.RPL_NONE:
		// Undefined format
		break

	case utils.RPL_WHOISREGNICK:
		// <client> <nick> :has identified for this nick
		text := strings.Join(message.Parameters[1:], " ")
		s.BufferWho = append(s.BufferWho, text)

	case utils.RPL_WHOISUSER:
		// <client> <nick> <username> <host> * :<realname>
		text := strings.Join(message.Parameters[1:], " ")
		s.BufferWho = append(s.BufferWho, text)

	case utils.RPL_WHOISSERVER:
		// <client> <nick> <server> :<server info>
		text := strings.Join(message.Parameters[1:], " ")
		s.BufferWho = append(s.BufferWho, text)

	case utils.RPL_WHOISOPERATOR:
		// <client> <nick> :is an IRC operator
		text := strings.Join(message.Parameters[1:], " ")
		s.BufferWho = append(s.BufferWho, text)

	case utils.RPL_WHOWASUSER:
		// <client> <nick> <username> <host> * :<realname>
		text := strings.Join(message.Parameters[1:], " ")
		s.BufferWho = append(s.BufferWho, text)

	case utils.RPL_ENDOFWHO:
		// <client> <mask> :End of WHO list
		for _, m := range s.BufferWho {
			s.log(fmt.Sprintf("[who] %s", m))
		}
		s.BufferWho = make([]string, 0)

	case utils.RPL_WHOISIDLE:
		// <client> <nick> <secs> <signon> :seconds idle, signon time
		text := strings.Join(message.Parameters[1:], " ")
		s.BufferWho = append(s.BufferWho, text)

	case utils.RPL_ENDOFWHOIS:
		// <client> <nick> :End of /WHOIS list
		for _, m := range s.BufferWho {
			s.log(fmt.Sprintf("[whois] %s", m))
		}
		s.BufferWho = make([]string, 0)

	case utils.RPL_WHOISCHANNELS:
		// <client> <nick> :[prefix]<channel>{ [prefix]<channel>}
		text := fmt.Sprintf("%s is in channels %s", message.Parameters[1], message.Parameters[2])
		s.BufferWho = append(s.BufferWho, text)

	case utils.RPL_WHOISSPECIAL:
		// <client> <nick> :blah blah blah
		text := strings.Join(message.Parameters[1:], " ")
		s.BufferWho = append(s.BufferWho, text)

	case utils.RPL_LISTSTART:
		// <client> <channel> :Users...
		s.bufferList = append(s.bufferList, message.Parameters[2])

	case utils.RPL_LIST:
		// <client> <channel> <client count> :<topic>
		s.bufferList = append(s.bufferList, message.Parameters[3])

	case utils.RPL_LISTEND:
		// <client> :End of /LIST
		for _, m := range s.bufferList {
			s.log(fmt.Sprintf("[list] %s", m))
		}
		s.bufferList = make([]string, 0)

	case utils.RPL_CREATIONTIME:
		// <client> <channel> <creationtime>
		timestamp, _ := strconv.ParseInt(message.Parameters[2], 10, 64)
		date := time.Unix(timestamp, 0).String()
		s.log(fmt.Sprintf("%s was created on %s", message.Parameters[1], date))

	case utils.RPL_WHOISACCOUNT:
		// <client> <nick> <account> :is logged in as
		text := fmt.Sprintf("%s %s %s", message.Parameters[1], message.Parameters[3], message.Parameters[2])
		s.BufferWho = append(s.BufferWho, text)

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

	case utils.RPL_WHOISACTUALLY:
		// <client> <nick> [<host> [<ip>]] :Is actually...
		paramCount := len(message.Parameters)
		text := fmt.Sprintf("%s %s", message.Parameters[1:paramCount-1], message.Parameters[paramCount-1])
		s.BufferWho = append(s.BufferWho, text)

	case utils.RPL_VERSION:
		// <client> <version> <server> :<comments>
		s.log(fmt.Sprintf("%s %s %s", message.Parameters[2], message.Parameters[1], message.Parameters[3]))

	case utils.RPL_WHOREPLY:
		// <client> <channel> <username> <host> <server> <nick> <flags> :<hopcount> <realname>
		text := strings.Join(message.Parameters[1:], " ")
		s.BufferWho = append(s.BufferWho, text)

	case utils.RPL_NAMREPLY:
		// <client> <symbol> <channel> :[prefix]<nick>{ [prefix]<nick>}
		s.channelsJoined[message.Parameters[2]].usersJoin(strings.Split(message.Parameters[3], " "))

	case utils.RPL_LINKS:
		// <client> * <server> :<hopcount> <server info>
		text := fmt.Sprintf("%s %s", message.Parameters[2], message.Parameters[3])
		s.bufferLinks = append(s.bufferHelp, text)

	case utils.RPL_ENDOFLINKS:
		// <client> * :End of /LINKS list
		for _, m := range s.bufferLinks {
			s.log(fmt.Sprintf("[links] %s", m))
		}
		s.bufferLinks = make([]string, 0)

	case utils.RPL_ENDOFNAMES:
		// <client> <channel> :End of /NAMES list
		break

	case utils.RPL_MOTDSTART:
		// <client> :- <server> Message of the day -
		s.bufferMotd = append(s.bufferMotd, message.Parameters[1])

	case utils.RPL_ENDOFWHOWAS:
		// <client> <nick> :End of WHOWAS
		for _, m := range s.BufferWho {
			s.log(fmt.Sprintf("[whowas] %s", m))
		}
		s.BufferWho = make([]string, 0)

	case utils.RPL_INFO:
		// <client> :<string>
		s.bufferInfo = append(s.bufferInfo, message.Parameters[1])

	case utils.RPL_MOTD:
		// <client> :<line of the motd>
		s.bufferMotd = append(s.bufferMotd, message.Parameters[1])

	case utils.RPL_ENDOFINFO:
		// <client> :End of INFO list
		for _, m := range s.bufferInfo {
			s.log(fmt.Sprintf("[info] %s", m))
		}
		s.bufferInfo = make([]string, 0)

	case utils.RPL_ENDOFMOTD:
		// <client> :End of /MOTD command.
		for _, m := range s.bufferMotd {
			s.log(fmt.Sprintf("[motd] %s", m))
		}
		s.bufferMotd = make([]string, 0)

	case utils.RPL_WHOISHOST:
		// <client> <nick> :is connecting from *@localhost 127.0.0.1
		text := fmt.Sprintf("%s %s", message.Parameters[1], message.Parameters[2])
		s.BufferWho = append(s.BufferWho, text)

	case utils.RPL_WHOISMODES:
		// <client> <nick> :is using modes +ailosw
		text := fmt.Sprintf("%s %s", message.Parameters[1], message.Parameters[2])
		s.BufferWho = append(s.BufferWho, text)

	case utils.RPL_TIME:
		// <client> <server> [<timestamp> [<TS offset>]] :<human-readable time>
		s.log(fmt.Sprintf("Time on %s is %s", message.Parameters[1], message.Parameters[len(message.Parameters)-1]))

	case utils.RPL_WHOISSECURE:
		// <client> <nick> :is using a secure connection
		text := fmt.Sprintf("%s %s", message.Parameters[1], message.Parameters[2])
		s.BufferWho = append(s.BufferWho, text)

	case utils.RPL_HELPSTART:
		// <client> <subject> :<first line of help section>
		s.bufferHelp = append(s.bufferHelp, message.Parameters[2])

	case utils.RPL_HELPTXT:
		// <client> <subject> :<line of help text>
		s.bufferHelp = append(s.bufferHelp, message.Parameters[2])

	case utils.RPL_ENDOFHELP:
		// <client> <subject> :<last line of help text>
		for _, b := range s.bufferHelp {
			s.log(fmt.Sprintf("[help] %s", b))
		}
		s.bufferHelp = make([]string, 0)

	case utils.ERR_UNKNOWNERROR,
		utils.ERR_NOSUCHNICK,
		utils.ERR_NOSUCHSERVER,
		utils.ERR_NOSUCHCHANNEL,
		utils.ERR_CANNOTSENDTOCHAN,
		utils.ERR_TOOMANYCHANNELS,
		utils.ERR_WASNOSUCHNICK,
		utils.ERR_NOORIGIN,
		utils.ERR_NORECIPIENT,
		utils.ERR_NOTEXTTOSEND,
		utils.ERR_INPUTTOOLONG,
		utils.ERR_UNKNOWNCOMMAND,
		utils.ERR_NOMOTD,
		utils.ERR_NONICKNAMEGIVEN,
		utils.ERR_ERRONEUSNICKNAME,
		utils.ERR_NICKNAMEINUSE,
		utils.ERR_NICKCOLLISION,
		utils.ERR_USERNOTINCHANNEL,
		utils.ERR_NOTONCHANNEL,
		utils.ERR_USERONCHANNEL,
		utils.ERR_NOTREGISTERED,
		utils.ERR_NEEDMOREPARAMS,
		utils.ERR_ALREADYREGISTERED,
		utils.ERR_PASSWDMISMATCH,
		utils.ERR_YOUREBANNEDCREEP,
		utils.ERR_CHANNELISFULL,
		utils.ERR_UNKNOWNMODE,
		utils.ERR_INVITEONLYCHAN,
		utils.ERR_BANNEDFROMCHAN,
		utils.ERR_BADCHANNELKEY,
		utils.ERR_BADCHANMASK,
		utils.ERR_NOPRIVILEGES,
		utils.ERR_CHANOPRIVSNEEDED,
		utils.ERR_CANTKILLSERVER,
		utils.ERR_NOOPERHOST,
		utils.ERR_UMODEUNKNOWNFLAG,
		utils.ERR_USERSDONTMATCH,
		utils.ERR_HELPNOTFOUND,
		utils.ERR_INVALIDKEY,
		utils.ERR_STARTTLS,
		utils.ERR_INVALIDMODEPARAM,
		utils.ERR_NOPRIVS,
		utils.ERR_NICKLOCKED,
		utils.ERR_SASLFAIL,
		utils.ERR_SASLTOOLONG,
		utils.ERR_SASLABORTED,
		utils.ERR_SASLALREADY:
		paramCount := len(message.Parameters)
		text := message.Parameters[paramCount-1]
		if paramCount > 1 {
			text += fmt.Sprintf(" (%s)", strings.Join(message.Parameters[1:paramCount-1], ", "))
		}
		s.logs.Append(s.host, utils.LogError, text)
		s.BufferWho = make([]string, 0)
		s.BufferStats = make([]string, 0)

	default:
		text := fmt.Sprintf("Unimplemented reply: %s", utils.MarshalMessage(message))
		s.logs.Append("System", utils.LogError, text)
	}
}
