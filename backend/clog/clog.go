package clog

import (
	"log"
)

// TODO: Debug logging
func Debug(v ...any) {
	log.Println(v...)
}

func Debugf(format string, v ...any) {
	log.Printf(format, v...)
}

func Println(v ...any) {
	log.Println(v...)
}

func Printf(format string, v ...any) {
	log.Printf(format, v...)
}

func Fatal(v ...any) {
	log.Fatal(v...)
}
