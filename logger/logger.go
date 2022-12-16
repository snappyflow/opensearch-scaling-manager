package logger

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"gopkg.in/natefinch/lumberjack.v2"
)

// loggers

var (
	Trace *log.Logger
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
	Debug *log.Logger
	Fatal *log.Logger
)

type LOG struct {
	module string
	Trace  *log.Logger
	Info   *log.Logger
	Warn   *log.Logger
	Error  *log.Logger
	Debug  *log.Logger
	Fatal  *log.Logger
}

var k = koanf.New(".")

// Init initilise logging
func (l *LOG) Init(module string) {
	l.module = module

	var (
		TRACE   = fmt.Sprintf("%5s%15s ", "TRACE", module)
		INFO    = fmt.Sprintf("%5s%15s ", "INFO", module)
		WARNING = fmt.Sprintf("%5s%15s ", "WARN", module)
		ERROR   = fmt.Sprintf("%5s%15s ", "ERROR", module)
		DEBUG   = fmt.Sprintf("%5s%15s ", "DEBUG", module)
		FATAL   = fmt.Sprintf("%5s%15s ", "FATAL", module)
	)

	if err := k.Load(file.Provider("logger/log_config.json"), json.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	if _, err := os.Stat(k.String("logpath")); os.IsNotExist(err) {
		// Path does not exist, create necessary folders in specified path
		err := os.MkdirAll(k.String("logpath"), os.ModePerm)
		if err != nil {
			log.Fatalf("Error creating dir: %v", err)
		}
	}

	// create log path
	path := path.Join(k.String("logpath"), k.String("logfile"))

	// create lumberjack loger object for rotaing file handling
	logger := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    k.Int("MaxSize"), // megabytes
		MaxBackups: k.Int("MaxBackups"),
		MaxAge:     k.Int("MaxAge"), // days
	}
	// get log level and convert to upper case for switch statement
	level := strings.ToUpper(k.String("level"))

	traceHandle := io.MultiWriter(ioutil.Discard, logger)
	infoHandle := io.MultiWriter(os.Stdout, logger)
	warningHandle := io.MultiWriter(os.Stdout, logger)
	errorHandle := io.MultiWriter(os.Stderr, logger)
	debugHandle := io.MultiWriter(os.Stdout, logger)
	fatalHandle := io.MultiWriter(os.Stdout, logger)

	l.Trace = log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile)
	l.Debug = log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile)
	l.Info = log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile)
	l.Warn = log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile)
	l.Error = log.New(ioutil.Discard, ERROR, log.Ldate|log.Ltime|log.Lshortfile)
	l.Fatal = log.New(ioutil.Discard, FATAL, log.Ldate|log.Ltime|log.Lshortfile)

	//setup loging based on level
	switch level {
	default:
		l.Trace = log.New(traceHandle, TRACE, log.Ldate|log.Ltime|log.Lshortfile)
		l.Debug = log.New(debugHandle, DEBUG, log.Ldate|log.Ltime|log.Lshortfile)
		l.Info = log.New(infoHandle, INFO, log.Ldate|log.Ltime|log.Lshortfile)
		l.Warn = log.New(warningHandle, WARNING, log.Ldate|log.Ltime|log.Lshortfile)
		l.Error = log.New(errorHandle, ERROR, log.Ldate|log.Ltime|log.Lshortfile)
		l.Fatal = log.New(fatalHandle, FATAL, log.Ldate|log.Ltime|log.Lshortfile)
	case "TRACE":
		l.Trace = log.New(traceHandle, TRACE, log.Ldate|log.Ltime|log.Lshortfile)
		fallthrough
	case "DEBUG":
		l.Debug = log.New(debugHandle, DEBUG, log.Ldate|log.Ltime|log.Lshortfile)
		fallthrough
	case "INFO":
		l.Info = log.New(infoHandle, INFO, log.Ldate|log.Ltime|log.Lshortfile)
		fallthrough
	case "WARNING":
		l.Warn = log.New(warningHandle, WARNING, log.Ldate|log.Ltime|log.Lshortfile)
		fallthrough
	case "ERROR":
		l.Error = log.New(errorHandle, ERROR, log.Ldate|log.Ltime|log.Lshortfile)
	case "FATAL":
		l.Fatal = log.New(fatalHandle, FATAL, log.Ldate|log.Ltime|log.Lshortfile)
	}
}
