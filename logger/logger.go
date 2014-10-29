package logger

import (
	"io"
	"log"
	"os"
)

//logger
var mylogger *log.Logger = nil
var logf *os.File = nil

func FailOnError(err error, msg string) {
	if err != nil && mylogger != nil {
		mylogger.Println(msg, err)
	}
}
func CreateLogger(logFileName string) {
	var err error
	logf, err = os.OpenFile(logFileName, os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		logf, err = os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE, 0640)
		if err != nil {
			log.Fatalln(err)
		}
	}

	log.SetOutput(logf)
	mylogger = log.New(io.MultiWriter(logf, os.Stdout), "worker: ", log.Ldate|log.Ltime|log.Llongfile)
	mylogger.Println("goes to logf")
}
func DropLogger() {
	if logf != nil {
		logf.Close()
	}
}
