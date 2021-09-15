package mock

import "errors"

// Implemented Logger interface, for use in testing. Doesn't log instead storing errors and messages. Also contains some utility functions
type MockLogger struct {
	Errors   []error
	Debugs   []error
	Warnings []string
	Infos    []string
}

func (l *MockLogger) Error(err error) {
	l.Errors = append(l.Errors, err)
}

func (l *MockLogger) Warning(msg string) {
	l.Warnings = append(l.Warnings, msg)
}

func (l *MockLogger) Info(msg string) {
	l.Infos = append(l.Infos, msg)
}

func (l *MockLogger) Debug(err error) {
	l.Debugs = append(l.Debugs, err)
}

func (l *MockLogger) InWarnings(str string) bool {
	return inStrArr(str, l.Warnings)
}

func (l *MockLogger) InInfos(str string) bool {
	return inStrArr(str, l.Infos)
}

func (l *MockLogger) InErrors(err error) bool {
	return inErrArr(err, l.Errors)
}

func (l *MockLogger) InDebugs(err error) bool {
	return inErrArr(err, l.Debugs)
}

func (l *MockLogger) NotInWarnings(str string) bool {
	return !inStrArr(str, l.Warnings)
}

func (l *MockLogger) NotInInfos(str string) bool {
	return !inStrArr(str, l.Infos)
}

func (l *MockLogger) NotInErrors(err error) bool {
	return !inErrArr(err, l.Errors)
}

func (l *MockLogger) NotInDebugs(err error) bool {
	return !inErrArr(err, l.Debugs)
}

func (l *MockLogger) StrInErrors(str string) bool {
	return strInErrArr(str, l.Errors)
}

func (l *MockLogger) StrInDebugs(str string) bool {
	return strInErrArr(str, l.Debugs)
}

func (l *MockLogger) NotStrInErrors(str string) bool {
	return !strInErrArr(str, l.Errors)
}

func (l *MockLogger) NotStrInDebugs(str string) bool {
	return !strInErrArr(str, l.Debugs)
}

func inStrArr(str string, arr []string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}

func strInErrArr(str string, arr []error) bool {
	for _, e := range arr {
		if str == e.Error() {
			return true
		}
	}
	return false
}

func inErrArr(err error, arr []error) bool {
	for _, e := range arr {
		if errors.Is(e, err) {
			return true
		}
	}
	return false
}
