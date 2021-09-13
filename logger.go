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
	l.infoFunc = l.doNothing
	l.errFunc = l.doNothingErr
	l.warFunc = l.doNothing
	l.debFunc = l.doNothingErr
	if options.Contains(Empty) {
		// default is to do nothing
	} else if options.Contains(Quiet) {
		l.errLogger = log.New(os.Stderr, "ERROR: ", log.LstdFlags) // For fatal or big errors
		l.errFunc = l.error
	} else {
		l.infoFunc = l.info
		l.errFunc = l.error
		l.warFunc = l.warning
		l.debFunc = l.debug
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
	infoFunc   func(string)
	warFunc    func(string)
	errFunc    func(error)
	debFunc    func(error)
	doDebug    bool
	mut        sync.Mutex
}

func (l *logger) Info(msg string) {
	l.infoFunc(msg)
}

func (l *logger) Warning(msg string) {
	l.warFunc(msg)
}

func (l *logger) Error(err error) {
	l.errFunc(err)
}

func (l *logger) Debug(err error) {
	if l.doDebug {
		l.debFunc(err)
	}
}

func (l *logger) info(msg string) {
	l.infoLogger.Println(msg)
}

func (l *logger) error(err error) {
	l.errLogger.Println(err)
}

func (l *logger) warning(msg string) {
	l.warLogger.Println(msg)
}

func (l *logger) debug(err error) {
	if l.doDebug {
		l.debLogger.Println(err)
	}
}

func (l *logger) doNothing(string)   {}
func (l *logger) doNothingErr(error) {}
