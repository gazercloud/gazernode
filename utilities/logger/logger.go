package logger

import (
	"fmt"
	"github.com/gazercloud/gazernode/utilities"
	"log"
	"os"
	"path/filepath"
	"time"
)

var logsPath string
var loggerObject *log.Logger
var currentLogFileName string
var currentLogFile *os.File

var LogDepthDays = 30

func CurrentExePath() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}

func Init(path string) {
	logsPath = path
	err := os.MkdirAll(path, 0777)
	if err != nil {
		fmt.Println("Can not create log directory")
	}
}

func DefaultLogPath() string {
	return CurrentExePath() + "/logs/web"
}

func CheckLogFile() {
	var err error
	logFile := logsPath + "/" + time.Now().Format("2006-01-02") + ".log"

	if logFile != currentLogFileName {
		if currentLogFile != nil {
			_ = currentLogFile.Close()
		}

		if loggerObject != nil {
			loggerObject = nil
		}

		currentLogFile, err = os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("error opening file: ", err)
		}

		loggerObject = log.New(currentLogFile, "", log.Ldate|log.Lmicroseconds)
		time.Sleep(time.Millisecond * 500)
		currentLogFileName = logFile

		files, err := utilities.GetDir(logsPath)
		if err == nil {
			for _, file := range files {
				if !file.Dir {
					t, err := time.Parse("2006-01-02", file.NameWithoutExt)
					if err == nil {
						if time.Now().Sub(t) > time.Duration(LogDepthDays*24)*time.Hour {
							_ = os.Remove(file.Path)
						}
					}
				}
			}
		}
	}
}

func Println(v ...interface{}) {
	CheckLogFile()
	if loggerObject != nil {
		loggerObject.Println(v)
	}
	fmt.Print(time.Now().UTC().Format("2006-01-02 15:04:05.999"), " ")
	fmt.Println(v...)
}

func Error(v ...interface{}) {
	CheckLogFile()
	if loggerObject != nil {
		loggerObject.Println(v)
	}
	fmt.Print(time.Now().UTC().Format("2006-01-02 15:04:05.999"), " ")
	fmt.Println(v...)
}
