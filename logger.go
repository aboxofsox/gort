package gort

import (
	"fmt"
	"io"
	"runtime"
	"time"
)

type Level int

const (
	INFO Level = iota
	WARNING
	ERROR
	DEBUG
)

var Levels = map[Level]string{
	INFO:    "INFO",
	WARNING: "WARNING",
	ERROR:   "ERROR",
	DEBUG:   "DEBUG",
}

type Logger struct {
	w      io.Writer
	events chan Event
}

type Event struct {
	Timestamp string
	Level     Level
	Message   string
}

func (e Event) String() string {
	return fmt.Sprintf("[%s] %s %s", e.Timestamp, Levels[e.Level], e.Message)
}

func NewLogger(w io.Writer) *Logger {
	logger := &Logger{
		w:      w,
		events: make(chan Event),
	}

	go logger.run()

	return logger
}

func (l *Logger) run() {
	select {
	case event := <-l.events:
		l.w.Write([]byte(event.String()))
	}
}

func (l *Logger) Log(level Level, message string) {
	l.events <- Event{
		Timestamp: time.Now().Format(time.ANSIC),
		Level:     level,
		Message:   message + "\n",
	}
}

func (l *Logger) Info(message string) {
	l.Log(INFO, message)
}

func (l *Logger) Warning(message string) {
	l.Log(WARNING, message)
}

func (l *Logger) Error(message string) {
	l.Log(ERROR, message)
}

func (l *Logger) Debug(message string) {
	stack := getStacktrace()
	l.Log(DEBUG, message+"\n"+stack)
}

func getStacktrace() string {
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	return string(buf[0:n])
}
