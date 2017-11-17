package log

import (
	"os"
	"strings"
	"time"
	"log"
)

// interface
type LoggerInterface interface {
	Create() error
	Write(...string)
	WriteError(err error)
	Close()
}

type Logger struct {
	// write to file is enable
	Disabled bool

	// ex. ./logs/app.log
	FileName string

	// descriptor
	Source os.File

	DateFormat string
}

// opening or creating log's file with FileName
func (l *Logger) Create() error {
	var err error = nil

	if !l.Disabled {
		var flag int
		if _, err := os.Stat(l.FileName); os.IsNotExist(err) {
			flag = os.O_RDWR|os.O_CREATE
		} else {
			flag = os.O_APPEND|os.O_WRONLY
		}

		f, err := os.OpenFile(l.FileName, flag, 0755)

		if err != nil {
			return err
		}

		l.Source = *f
	}

	return err
}

// writing messages in log file
func (l *Logger) Write(messages ...string) {
	if !l.Disabled {
		date := time.Now()

		str := strings.Replace("{date} {text}\n", "{date}", date.Format(l.DateFormat), -1)

		var text string
		for _, message := range messages {
			text += message + " "
		}

		str = strings.Replace(str, "{text}", text, -1)
		l.Source.WriteString(str)
	} else {
		log.Println(messages)
	}
}

// writing messages in log from errors
func (l *Logger) WriteError(err error) {
	if l.Disabled {
		log.Println(err)
	} else {
		l.Write(err.Error())
	}
}

// close
func (l *Logger) Close() {
	l.Source.Close()
}

// constructor
func NewLogger(fileName string, disabled bool) (*Logger, error) {
	logger := &Logger{
		FileName: fileName,
		DateFormat: time.RFC3339,
		Disabled: disabled,
	}
	err := logger.Create()
	return logger, err
}
