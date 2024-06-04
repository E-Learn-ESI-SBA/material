package logs

import (
	"fmt"
	"log"
	"madaurus/dev/material/app/configs"
	"madaurus/dev/material/app/interfaces"
	"os"
	"path/filepath"
	"runtime"
)

type Level int

var (
	F *os.File

	DefaultPrefix      = ""
	DefaultCallerDepth = 2

	logger      *log.Logger
	InfoLogger  *log.Logger = log.New(os.Stdout, "INFO", log.LstdFlags)
	ErrorLogger *log.Logger = log.New(os.Stderr, "ERROR", log.LstdFlags)
	DebugLogger *log.Logger = log.New(os.Stdout, "DEBUG", log.LstdFlags)
	logPrefix               = ""
	levelFlags              = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func Setup(AppSettings *interfaces.App) {
	var err error
	filePath := getLogFilePath(AppSettings)
	fileName := getLogFileName(AppSettings)
	log.Printf("filePath: %s", filePath)
	F, err := configs.MustOpen(fileName, filePath)
	if err != nil {
		log.Fatalf("logging.Setup err: %v", err)
	}
	logger = log.New(F, DefaultPrefix, log.LstdFlags)
}

func Debug(v ...interface{}) {
	setPrefix(DEBUG)
	logger.Println(v)
	return

}

// output logs at info level
func Info(v ...interface{}) {
	setPrefix(INFO)
	InfoLogger.Println(v)
	return

}

// output logs at warn level
func Warn(v ...interface{}) {
	setPrefix(WARNING)
	logger.Println(v)
	return

}

// output logs at error level
func Error(v ...interface{}) {
	setPrefix(ERROR)
	ErrorLogger.Println(v)
	return
}

// output logs at fatal level
func Fatal(v ...interface{}) {
	setPrefix(FATAL)
	logger.Println(v)
}

// set the prefix of the log output
func setPrefix(level Level) {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}

	logger.SetPrefix(logPrefix)
}
