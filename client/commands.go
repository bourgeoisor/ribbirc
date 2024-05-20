package client

import (
	"fmt"
	"ribbirc/utils"
	"strconv"
	"strings"
	"time"
)

func (s *Server) handleCommand(input string, channel string) *utils.Message {
	message := &utils.Message{}

	parts := strings.Split(input, " ")
	paramCount := len(parts) - 1

	switch parts[0] {
	case "/ping":
		if paramCount > 0 {
			s.invalidCommandParameters("/ping")
			return nil
		}
		message.Command = "PING"
		message.Parameters = []string{strconv.FormatInt(time.Now().UnixMilli(), 10)}

	case "/quit":
		message.Command = "QUIT"
		if paramCount == 1 {
			message.Parameters = []string{strings.Join(parts[1:], " ")}
		}

	case "/nick":
		if paramCount != 1 {
			s.invalidCommandParameters("/nick <nickname>")
			return nil
		}
		message.Command = "NICK"
		message.Parameters = []string{parts[1]}

	case "/oper":
		if paramCount != 2 {
			s.invalidCommandParameters("/oper <name> <password>")
			return nil
		}
		message.Command = "OPER"
		message.Parameters = []string{parts[1], parts[2]}

	case "/motd":
		if paramCount > 1 {
			s.invalidCommandParameters("/motd [<target>]")
			return nil
		}
		message.Command = "MOTD"
		if paramCount == 1 {
			message.Parameters = []string{parts[1]}
		}

	case "/version":
		if paramCount > 1 {
			s.invalidCommandParameters("/version [<target>]")
			return nil
		}
		message.Command = "VERSION"
		if paramCount == 1 {
			message.Parameters = []string{parts[1]}
		}

	case "/admin":
		if paramCount > 1 {
			s.invalidCommandParameters("/admin [<target>]")
			return nil
		}
		message.Command = "ADMIN"
		if paramCount == 1 {
			message.Parameters = []string{parts[1]}
		}

	case "/connect":
		if paramCount > 3 {
			s.invalidCommandParameters("<target server> [<port> [<remote server>]]")
			return nil
		}
		message.Command = "CONNECT"
		message.Parameters = []string{parts[1]}
		if paramCount > 1 {
			message.Parameters = append(message.Parameters, parts[2])
		}
		if paramCount > 2 {
			message.Parameters = append(message.Parameters, parts[3])
		}

	case "/lusers":
		if paramCount > 0 {
			s.invalidCommandParameters("/lusers")
			return nil
		}
		message.Command = "LUSERS"

	case "/time":
		if paramCount > 1 {
			s.invalidCommandParameters("/time [<server>]")
			return nil
		}
		message.Command = "TIME"
		if paramCount == 1 {
			message.Parameters = []string{parts[1]}
		}

	case "/stats":
		if paramCount < 1 || paramCount > 2 {
			s.invalidCommandParameters("/stats <query> [<server>]")
			return nil
		}
		message.Command = "STATS"
		message.Parameters = []string{parts[1]}
		if paramCount == 2 {
			message.Parameters = append(message.Parameters, parts[2])
		}

	case "/help":
		if paramCount > 1 {
			s.invalidCommandParameters("/help [<subject>]")
			return nil
		}
		message.Command = "HELP"
		if paramCount == 1 {
			message.Parameters = []string{parts[1]}
		}

	case "/info":
		if paramCount > 0 {
			s.invalidCommandParameters("/info")
			return nil
		}
		message.Command = "INFO"

	case "/join":
		if paramCount < 1 || paramCount > 2 {
			s.invalidCommandParameters("/join <channel>{,<channel>} [<key>{,<key>}]")
			return nil
		}
		message.Command = "JOIN"
		message.Parameters = []string{parts[1]}
		if paramCount == 2 {
			message.Parameters = append(message.Parameters, parts[2])
		}

	case "/part":
		if paramCount < 1 {
			s.invalidCommandParameters("/part <channel>{,<channel>} [<reason>]")
			return nil
		}
		message.Command = "PART"
		message.Parameters = []string{parts[1]}
		if paramCount > 1 {
			message.Parameters = append(message.Parameters, strings.Join(parts[2:], " "))
		}

	case "/topic":
		if paramCount < 1 || paramCount > 2 {
			s.invalidCommandParameters("/topic <channel> [<topic>]")
			return nil
		}
		message.Command = "TOPIC"
		message.Parameters = []string{parts[1]}
		if paramCount > 1 {
			message.Parameters = append(message.Parameters, strings.Join(parts[2:], " "))
		}

	case "/names":
		if paramCount != 1 {
			s.invalidCommandParameters("/names <channel>{,<channel>}")
			return nil
		}
		message.Command = "NAMES"
		message.Parameters = []string{parts[1]}

	case "/list":
		if paramCount > 1 {
			s.invalidCommandParameters("/list [<channel>{,<channel>}] [<elistcond>{,<elistcond>}]")
			return nil
		}
		message.Command = "LIST"
		if paramCount == 1 {
			message.Parameters = []string{parts[1]}
		}

	case "/invite":
		if paramCount != 2 {
			s.invalidCommandParameters("/invite <nickname> <channel>")
			return nil
		}
		message.Command = "INVITE"
		message.Parameters = []string{parts[1], parts[2]}

	case "/kick":
		if paramCount < 2 {
			s.invalidCommandParameters("/kick <channel> <user> *( \",\" <user> ) [<comment>]")
			return nil
		}
		message.Command = "KICK"
		message.Parameters = []string{parts[1], parts[2]}
		if paramCount > 2 {
			message.Parameters = append(message.Parameters, strings.Join(parts[3:], " "))
		}

	case "/who":
		if paramCount != 1 {
			s.invalidCommandParameters("/who <mask>")
			return nil
		}
		message.Command = "WHO"
		message.Parameters = []string{parts[1]}

	case "/whois":
		if paramCount < 1 || paramCount > 2 {
			s.invalidCommandParameters("/whois [<target>] <nick>")
			return nil
		}
		message.Command = "WHOIS"
		message.Parameters = []string{parts[1]}
		if paramCount > 1 {
			message.Parameters = append(message.Parameters, parts[2])
		}

	case "/whowas":
		if paramCount < 1 || paramCount > 2 {
			s.invalidCommandParameters("/whowas <nick> [<count>]")
			return nil
		}
		message.Command = "WHOWAS"
		message.Parameters = []string{parts[1]}
		if paramCount > 1 {
			message.Parameters = append(message.Parameters, parts[2])
		}

	case "/kill":
		if paramCount < 2 {
			s.invalidCommandParameters("/kill <nickname> <comment>")
			return nil
		}
		message.Command = "KILL"
		message.Parameters = []string{parts[1], strings.Join(parts[2:], " ")}

	case "/rehash":
		if paramCount > 0 {
			s.invalidCommandParameters("/rehash")
			return nil
		}
		message.Command = "REHASH"

	case "/restart":
		if paramCount > 0 {
			s.invalidCommandParameters("/restart")
			return nil
		}
		message.Command = "RESTART"

	case "/squit":
		if paramCount < 2 {
			s.invalidCommandParameters("/squit <server> <comment>")
			return nil
		}
		message.Command = "SQUIT"
		message.Parameters = []string{parts[1], strings.Join(parts[2:], " ")}

	case "/away":
		message.Command = "AWAY"
		if paramCount > 0 {
			message.Parameters = []string{strings.Join(parts[1:], " ")}
		}

	case "/links":
		if paramCount > 0 {
			s.invalidCommandParameters("/links")
			return nil
		}
		message.Command = "LINKS"

	case "/userhost":
		if paramCount > 5 {
			s.invalidCommandParameters("/userhost <nickname>{ <nickname>}")
			return nil
		}
		message.Command = "USERHOST"
		message.Parameters = parts[1:]

	case "/wallops":
		message.Command = "WALLOPS"
		message.Parameters = []string{strings.Join(parts[1:], " ")}

	default:
		text := fmt.Sprintf("Unimplemented command %s", parts[0])
		s.logs.Append("System", utils.LogError, text)
	}

	return message
}

func (s *Server) invalidCommandParameters(format string) {
	text := fmt.Sprintf("Invalid command format, expected '%s'.", format)
	s.logs.Append("System", utils.LogError, text)
}
