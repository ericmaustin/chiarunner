package main

import (
	"io"
	"log"
	"os"
	"time"
)

type logWriter struct {
	writer io.Writer
}

func (w logWriter) Write(bytes []byte) (int, error) {
	return w.writer.Write(append([]byte(time.Now().UTC().Format(time.RFC3339)+" "), bytes...))
}

func initLogger(logFileStr *string) {
	var writer io.Writer

	if logFileStr != nil && len(*logFileStr) > 0 {
		logFile, err := os.OpenFile(*logFileStr, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}
		writer = io.MultiWriter(os.Stdout, logFile)
	} else {
		writer = os.Stdout
	}

	log.SetFlags(0)
	log.SetOutput(&logWriter{
		writer: writer,
	})
}

func logF(fmt string, v ...interface{}) {
	log.Printf("[INFO] "+fmt, v...)
}

func logLn(v ...interface{}) {
	log.Println(append([]interface{}{"[INFO]"}, v...)...)
}

func logErrF(fmt string, v ...interface{}) {
	log.Printf("[ERROR] "+fmt, v...)
}

func logErrLn(v ...interface{}) {
	log.Println(append([]interface{}{"[ERROR]"}, v...)...)
}

func logFatalF(fmt string, v ...interface{}) {
	log.Fatalf("[FATAL] " + fmt, v...)
}

func logFatalLn(v ...interface{}) {
	log.Fatalln(append([]interface{}{"[FATAL]"}, v...)...)
}

