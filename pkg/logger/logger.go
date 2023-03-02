/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
	c.Println(fmt.Sprintf("\r ✅ "+msg, args...))
}

func (l *Logger) ImageRedLine(msg string, args ...interface{}) {
	c := color.New(color.FgHiRed)
	c.Println(fmt.Sprintf("\r ❌ "+msg, args...))
}
