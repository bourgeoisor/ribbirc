package utils

import "sync"

type Log struct {
	Source string
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

func (l *Logger) Append(source string, text string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.logs = append(l.logs, Log{source, text})
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
