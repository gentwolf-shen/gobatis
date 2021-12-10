package gobatis

import (
	"fmt"
	"log"
	"os"
)

var (
	logger ILogger
)

// 设置默认logger
func SetDefaultLogger() {
	o := &SystemLogger{}
	o.log = log.New(os.Stdout, "[gobatis] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	logger = o
}

// 设置默认第三方log
func SetCustomLogger(l ILogger) {
	logger = l
}

// 第三方logger需要实现的接口
type ILogger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
}

type SystemLogger struct {
	log *log.Logger
}

func (l *SystemLogger) output(prefix string, v ...interface{}) {
	_ = l.log.Output(3, "["+prefix+"] "+fmt.Sprint(v...))
}

func (l *SystemLogger) Debug(v ...interface{}) {
	l.output("DEBUG", v...)
}

func (l *SystemLogger) Info(v ...interface{}) {
	l.output("INFO", v...)
}

func (l *SystemLogger) Warn(v ...interface{}) {
	l.output("WARN", v...)
}

func (l *SystemLogger) Error(v ...interface{}) {
	l.output("ERROR", v...)
}
