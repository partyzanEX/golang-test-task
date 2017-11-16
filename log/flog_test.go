package log

import (
	"testing"
	"strings"
	"os"
	"errors"
)

var fileLog = "./app.log"

func TestNewLogger(t *testing.T) {
	logger, err := NewLogger(fileLog, false)
	checkError("start", err, t)
	logger.Write("Started log")
	logger.Write("Closed log")
	logger.Close()

	defer func() {
		if _, err := os.Stat(fileLog); !os.IsNotExist(err) {
			checkError("check log: exists", err, t)
		}

		file, err := os.Open(fileLog)
		checkError("check log: open", err, t)
		defer file.Close()
		defer os.Remove(fileLog)

		stat, err := file.Stat()
		if err != nil {
			t.Error(err)
		}

		buff := make([]byte, stat.Size())
		_, err = file.Read(buff)
		checkError("check log: read", err, t)

		result := string(buff)

		if !strings.Contains(result, "Started log") {
			t.Error("Don't started log")
		}

		if !strings.Contains(result, "Closed log") {
			t.Error("Don't closed log")
		}
	}()
}

func TestLogger_Write(t *testing.T) {
	logger, err := NewLogger(fileLog, false)
	checkError("write: start", err, t)
	defer logger.Close()

	test := "test message"
	logger.Write(test)

	defer func() {
		file, err := os.Open(fileLog)
		checkError("write: open", err, t)
		defer file.Close()
		defer os.Remove(fileLog)

		stat, err := file.Stat()
		if err != nil {
			t.Error(err)
		}

		buff := make([]byte, stat.Size())
		_, err = file.Read(buff)
		checkError("write: read", err, t)

		result := string(buff)

		if !strings.Contains(result, test) {
			t.Error("Don't write test")
		}
	}()
}

func TestLogger_WriteError(t *testing.T) {
	logger, err := NewLogger(fileLog, false)
	checkError("write_error: start", err, t)
	defer logger.Close()

	test := "test error"
	testErr := errors.New(test)
	logger.WriteError(testErr)

	defer func() {
		file, err := os.Open(fileLog)
		checkError("write_error: open", err, t)
		defer file.Close()
		defer os.Remove(fileLog)

		stat, err := file.Stat()
		if err != nil {
			t.Error(err)
		}

		buff := make([]byte, stat.Size())
		_, err = file.Read(buff)
		checkError("write_error: read", err, t)

		result := string(buff)

		if !strings.Contains(result, test) {
			t.Error("Don't write test error")
		}
	}()
}

func checkError(label string, err error, t *testing.T)  {
	if err != nil {
		t.Error(label)
		t.Error(err)
	}
}
