package logger

import (
	"log"
	"os"
	"sync"
)

type LoggerOptions uint32

func (o LoggerOptions) Contains(opt LoggerOptions) bool {
	return o&opt != 0
}

const (
	Debug LoggerOptions = 1 << iota
	Quiet               // Only show fatal/important error messages
	Empty               // don't print anything
)

type Logger interface {
	Info(string)
	Warning(string)
	Error(error)
	Debug(error)
}

func NewLogger(options LoggerOptions) *logger {
	l := new(logger)
	if options.Contains(Debug) {
		l.doDebug = true
	}
	l.infoFunc = doNothing
	l.errFunc = doNothingErr
	l.warFunc = doNothing
	l.debFunc = doNothingErr
	if options.Contains(Empty) {
		// default is to do nothing
	} else if options.Contains(Quiet) {
		l.errLogger = log.New(os.Stderr, "ERROR: ", log.LstdFlags) // For fatal or big errors
		l.errFunc = logErr
	} else {
		l.infoFunc = logStr
		l.errFunc = logErr
		l.warFunc = logStr
		l.debFunc = logErr
		l.infoLogger = log.New(os.Stderr, "INFO: ", log.LstdFlags)                      // inform the developer
		l.errLogger = log.New(os.Stderr, "ERROR: ", log.LstdFlags)                      // For fatal or big errors
		l.warLogger = log.New(os.Stderr, "WARNING: ", log.LstdFlags)                    //non-fatal errors
		l.debLogger = log.New(os.Stderr, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile) // more errors
	}
	return l
}

type logger struct {
	errLogger  *log.Logger
	infoLogger *log.Logger
	warLogger  *log.Logger
	debLogger  *log.Logger
	infoFunc   func(*log.Logger, string)
	warFunc    func(*log.Logger, string)
	errFunc    func(*log.Logger, error)
	debFunc    func(*log.Logger, error)
	doDebug    bool
	mut        sync.Mutex
}

func (l *logger) Info(msg string) {
	l.infoFunc(l.infoLogger, msg)
}

func (l *logger) Warning(msg string) {
	l.warFunc(l.warLogger, msg)
}

func (l *logger) Error(err error) {
	l.errFunc(l.errLogger, err)
}

func (l *logger) Debug(err error) {
	if l.doDebug {
		l.debFunc(l.debLogger, err)
	}
}

func logStr(l *log.Logger, msg string) {
	l.Println(msg)
}

func logErr(l *log.Logger, err error) {
	l.Println(err)
}

func doNothing(*log.Logger, string)   {}
func doNothingErr(*log.Logger, error) {}
