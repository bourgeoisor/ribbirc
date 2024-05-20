package utils

import "sync"

type LogKind int

const (
	LogPrivMsg LogKind = iota
	LogSystem
	LogError
	LogStatus
	LogJoined
	LogLeft
)

type Log struct {
	Source string
	Kind   LogKind
	Text   string
}

type Logger struct {
	mutex  sync.Mutex
	logs   []Log
	length int
}

func NewLogger() *Logger {
	return &Logger{
		logs: make([]Log, 0),
	}
}

func (l *Logger) Append(source string, kind LogKind, text string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.logs = append(l.logs, Log{source, kind, text})
	l.length++
}

func (l *Logger) GetNLogs(height int, offset int) []Log {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	start := l.length - height - offset
	end := l.length - offset

	logs := make([]Log, 0)
	for i := start; i < end; i++ {
		if i >= 0 {
			logs = append(logs, l.logs[i])
		}
	}
	return logs
}

func (l *Logger) GetAllLogs() []Log {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	return append([]Log{}, l.logs...)
}
