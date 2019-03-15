package logger

import (
	"log"
	"os"
)

type FileLogger struct {
	Logger      *log.Logger
	LogFilePath string
	LogFile     *os.File
}

type Logger interface {
	Println(v ...interface{})
	Printf(fmt string, v ...interface{})

	Fatalln(v ...interface{})
	Fatalf(fmtStr string, v ...interface{})

	Close() error
}

func NewFileLogger(filePath string, prefix string, flags int) (*FileLogger, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	return &FileLogger{
		Logger:  log.New(file, prefix, flags),
		LogFile: file,

		LogFilePath: filePath,
	}, nil
}

func (l *FileLogger) Println(v ...interface{}) {
	l.Logger.Println(v...)
}

func (l *FileLogger) Printf(fmt string, v ...interface{}) {
	l.Logger.Printf(fmt, v...)
}

func (l *FileLogger) Fatalln(v ...interface{}) {
	l.Logger.Fatalln(v...)
}

func (l *FileLogger) Fatalf(fmtStr string, v ...interface{}) {
	l.Fatalf(fmtStr, v...)
}

func (l *FileLogger) Close() error {
	return l.LogFile.Close()
}
