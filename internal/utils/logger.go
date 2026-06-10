package utils

import (
	"fmt"
	"os"
	"time"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	gray    = "\033[90m"
	blue    = "\033[94m"
	green   = "\033[92m"
	yellow  = "\033[93m"
	red     = "\033[91m"
	magenta = "\033[95m"
)

func log(level Level, format string, args ...any) {
	ts := time.Now().Format("2006-01-02 15:04:05.000")

	var (
		color string
		tag   string
	)

	switch level {
	case DEBUG:
		color = blue
		tag = "DEBUG"
	case INFO:
		color = green
		tag = "INFO"
	case WARN:
		color = yellow
		tag = "WARN"
	case ERROR:
		color = red
		tag = "ERROR"
	case FATAL:
		color = magenta
		tag = "FATAL"
	}

	msg := fmt.Sprintf(format, args...)

	fmt.Printf(
		"%s%s[%s]%s %s%s%s %s\n",
		gray,
		ts,
		tag,
		reset,
		bold,
		color,
		msg,
		reset,
	)

	if level == FATAL {
		os.Exit(1)
	}
}

func Debug(format string, args ...any) {
	log(DEBUG, format, args...)
}

func Info(format string, args ...any) {
	log(INFO, format, args...)
}

func Warn(format string, args ...any) {
	log(WARN, format, args...)
}

func Error(format string, args ...any) {
	log(ERROR, format, args...)
}

func Fatal(format string, args ...any) {
	log(FATAL, format, args...)
}

func HandlExp(e error) {
	if e == nil {
		return
	}

	Error("%v", e)
	os.Exit(1)
}
