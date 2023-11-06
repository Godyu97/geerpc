package logger

import (
	"go.uber.org/zap"
	"fmt"
)

type GeeLogger interface {
	Info(a ...any)
	Error(a ...any)
}

var geeLogger GeeLogger = &defaultLog{}

func InitLogger(logger GeeLogger) {
	//logger init
	geeLogger = logger
}

func Info(a ...any) {
	geeLogger.Info(a...)
}

func Error(a ...any) {
	geeLogger.Error(a...)
}

type defaultLog struct{}

func (l *defaultLog) Info(a ...any) {
	zap.L().Info(fmt.Sprintln(a...))
}

func (l *defaultLog) Error(a ...any) {
	zap.L().Error(fmt.Sprintln(a...))
}
