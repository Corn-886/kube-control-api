package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

type LogrusFormatter struct{}

// 格式化loggrus输出

func (s *LogrusFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006/01/02 - 15:04:05.000")
	var callerName string
	if len(entry.Caller.Function) >= 40 {
		callerName = entry.Caller.Function[33:len(entry.Caller.Function)]
	} else {
		callerName = entry.Caller.Function
	}
	msg := fmt.Sprintf("[LOG] %s   | %-60v %3v | %7v | %s\n", timestamp, callerName, entry.Caller.Line, entry.Level.String(), entry.Message)
	return []byte(msg), nil
}
