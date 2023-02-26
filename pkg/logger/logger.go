package logger

import (
	"fmt"

	"github.com/fatih/color"
)

type Logger struct {
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	if msg == "" {
		fmt.Println("")
		return
	}

	fmt.Println(fmt.Sprintf(msg, args...))
}

func (l *Logger) Error(err error) {
	c := color.New(color.FgHiRed)
	c.Println(fmt.Sprintf("%#v", err))
}

func (l *Logger) Header(msg string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(msg, args...))
}

func (l *Logger) StartImageLine(msg string, args ...interface{}) {
	c := color.New(color.FgHiYellow)
	c.Printf(fmt.Sprintf(msg, args...))
}

func (l *Logger) ImageGreenLine(msg string, args ...interface{}) {
	c := color.New(color.FgHiGreen)
	c.Println(fmt.Sprintf("\r"+msg, args...))
}

func (l *Logger) ImageRedLine(msg string, args ...interface{}) {
	c := color.New(color.FgHiRed)
	c.Println(fmt.Sprintf("\r"+msg, args...))
}
