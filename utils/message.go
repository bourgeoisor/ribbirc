package utils

import (
	"fmt"
	"strings"
)

type Message struct {
	tags       string
	Source     string
	Command    string
	Parameters []string
}

func (m *Message) SourceNick() string {
	index := strings.IndexByte(m.Source, '!')
	if index == -1 {
		return ""
	}
	return m.Source[:index]
}

func UnmarshalMessage(payload string) *Message {
	message := Message{}

	parts := strings.Split(payload, " ")
	for n, part := range parts {
		if message.Command == "" {
			switch part[0] {
			case '@':
				message.tags = part[1:]
			case ':':
				message.Source = part[1:]
			default:
				message.Command = part
			}
			continue
		}

		if part[0] != ':' {
			message.Parameters = append(message.Parameters, part)
		} else {
			message.Parameters = append(message.Parameters, strings.Join(parts[n:], " ")[1:])
			break
		}
	}

	return &message
}

func MarshalMessage(message *Message) string {
	var data string

	if message.tags != "" {
		data += fmt.Sprintf("@%s ", message.tags)
	}

	if message.Source != "" {
		data += fmt.Sprintf(":%s ", message.Source)
	}

	data += message.Command

	for i, parameter := range message.Parameters {
		if i == len(message.Parameters)-1 && strings.ContainsRune(parameter, ' ') {
			data += " :" + parameter
		} else {
			data += " " + parameter
		}
	}

	return data
}
