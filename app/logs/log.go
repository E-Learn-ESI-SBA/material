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

	logger     *log.Logger
	logPrefix  = ""
	levelFlags = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func Setup(AppSetting *interfaces.App) {
	var err error
	filePath := getLogFilePath(AppSetting)
	fileName := getLogFileName(AppSetting)
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
}

// output logs at info level
func Info(v ...interface{}) {
	setPrefix(INFO)
	logger.Println(v)
}

// output logs at warn level
func Warn(v ...interface{}) {
	setPrefix(WARNING)
	logger.Println(v)
}

// output logs at error level
func Error(v ...interface{}) {
	setPrefix(ERROR)
	logger.Println(v)
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
