package mock

import (
	"errors"
	"testing"

	"github.com/alexander-localbitcoins/logger"
)

type msgEnum uint32

const (
	InfoMsg msgEnum = 1 << iota
	WarMsg
	ErrMsg
	DebMsg
)

var msgs = map[msgEnum]string{
	InfoMsg: "SOMETHING INTERESTING HAPPENED",
	WarMsg:  "YOU MAY WANT TO CONSIDER THAT THIS HAPPENED",
	ErrMsg:  "THERE WAS AN ERROR",
	DebMsg:  "An extra error for more information for developers",
}

func sendToAll(l logger.Logger) (debErr error, errErr error) {
	l.Warning(msgs[WarMsg])
	l.Info(msgs[InfoMsg])
	debErr = errors.New(msgs[DebMsg])
	l.Debug(debErr)
	errErr = errors.New(msgs[ErrMsg])
	l.Error(errErr)
	return
}

func TestMockLogger(t *testing.T) {
	l := new(MockLogger)
	debErr, errErr := sendToAll(l)
	for k, m := range msgs {
		switch k {
		case InfoMsg:
			if len(l.Infos) != 1 || l.Infos[0] != m {
				t.Error("MockLogger infos does not contain send message")
			}
			if !l.InInfos(m) {
				t.Error("MockLogger InInfos did not find send message")
				continue
			}
			if l.NotInInfos(m) {
				t.Error("Inverse InInfos (NotInInfos) failed")
				continue
			}
		case WarMsg:
			if len(l.Warnings) != 1 || l.Warnings[0] != m {
				t.Error("MockLogger warnings does not contain send message")
			}
			if !l.InWarnings(m) {
				t.Error("MockLogger InWarnings did not find send message")
				continue
			}
			if l.NotInWarnings(m) {
				t.Error("Inverse InWarnings (NotInWarnings) failed")
				continue
			}
		case ErrMsg:
			if len(l.Errors) != 1 || !errors.Is(l.Errors[0], errErr) || l.Errors[0].Error() != m {
				t.Error("MockLogger errors does not contain send error")
			}
			if !l.InErrors(errErr) {
				t.Error("MockLogger InErrors did not find send error")
				continue
			}
			if l.NotInErrors(errErr) {
				t.Error("Inverse InErrors (NotInErrors) failed")
				continue
			}
			if !l.StrInErrors(m) {
				t.Error("Did not find error by error message")
				continue
			}
			if l.NotStrInErrors(m) {
				t.Error("Inverse of StrInError failed")
				continue
			}
		case DebMsg:
			if len(l.Debugs) != 1 || !errors.Is(l.Debugs[0], debErr) || l.Debugs[0].Error() != m {
				t.Error("MockLogger debugs does not contain send error")
			}
			if !l.InDebugs(debErr) {
				t.Error("MockLogger debugs did not contain send debug error")
				continue
			}
			if l.NotInDebugs(debErr) {
				t.Error("Inverse InDebugs (NotInDebugs) failed")
				continue
			}
			if !l.StrInDebugs(m) {
				t.Error("Did not find debug error by error message")
				continue
			}
			if l.NotStrInDebugs(m) {
				t.Error("Inverse of StrInDebugs failed")
				continue
			}
		}
	}
}
