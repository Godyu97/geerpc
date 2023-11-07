package logger

import (
	"fmt"
	"log"
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
	log.Print("INFO ", fmt.Sprintln(a...))
}

func (l *defaultLog) Error(a ...any) {
	log.Print("ERROR ", fmt.Sprintln(a...))
}
