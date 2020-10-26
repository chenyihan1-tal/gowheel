package logs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/xavierror/gowheel/file"
)

type (
	Level int
)

var (
	F *os.File

	DefaultCallerDepth = 2

	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

// Debug output logs at debug level
func Debug(v ...interface{}) {
	handle(DEBUG, v)
}

// DebugF output logs at debug level
func DebugF(format string, v ...interface{}) {
	handle(DEBUG, fmt.Sprintf(format, v))
}

// Info output logs at info level
func Info(v ...interface{}) {
	handle(INFO, v)
}

// InfoF output logs at info level
func InfoF(format string, v ...interface{}) {
	handle(INFO, fmt.Sprintf(format, v))
}

// Warn output logs at warn level
func Warn(v ...interface{}) {
	handle(WARNING, v)
}

// WarnF output logs at warn level
func WarnF(format string, v ...interface{}) {
	handle(WARNING, fmt.Sprintf(format, v))
}

// Error output logs at error level
func Error(v ...interface{}) {
	handle(ERROR, v)
}

// ErrorF output logs at error level
func ErrorF(format string, v ...interface{}) {
	handle(ERROR, fmt.Sprintf(format, v))
}

// Fatal output logs at fatal level
func Fatal(v ...interface{}) {
	handle(FATAL, v)
	os.Exit(0)
}

// FatalF output logs at fatal level
func FatalF(format string, v ...interface{}) {
	handle(FATAL, fmt.Sprintf(format, v))
	os.Exit(0)
}

func handle(level Level, v interface{}) {
	var err error
	var filePath string
	var fileName string
	var logContent string

	filePath = "runtime/logs/"

	fileName = levelFlags[level]
	filePath = filePath + time.Now().Format("/2006/0102/")

	F, err = file.MustOpen(fileName, filePath)
	if err != nil {
		log.Fatal()
	}

	timeNow := time.Now().Format("2006/01/02 15:04:05.000")

	// format logs
	_, f, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logContent = fmt.Sprintf("[%s][%s][%s:%d]%v", timeNow, levelFlags[level][0:1], filepath.Base(f), line, v)
	} else {
		logContent = fmt.Sprintf("[%s][%s]%v", timeNow, levelFlags[level][0:1], v)
	}

	fmt.Println(logContent)

	fmt.Fprintln(F, logContent)
}
