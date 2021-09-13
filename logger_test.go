package logger

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

type msgEnum uint32

const (
	InfoMsg msgEnum = 1 << iota
	WarMsg
	ErrMsg
	DebMsg
)

func (me msgEnum) Name() string {
	switch me {
	case InfoMsg:
		return "INFO"
	case WarMsg:
		return "WARNING"
	case ErrMsg:
		return "ERROR"
	case DebMsg:
		return "DEBUG"
	}
	return ""
}

var msgs = map[msgEnum]string{
	InfoMsg: "SOMETHING INTERESTING HAPPENED",
	WarMsg:  "YOU MAY WANT TO CONSIDER THAT THIS HAPPENED",
	ErrMsg:  "THERE WAS AN ERROR",
	DebMsg:  "An extra error for more information for developers",
}

type StdChanger struct {
	oldStdErr    *os.File
	newStdErr    *os.File
	stdErrOutput []byte
}

func (sc *StdChanger) Redirect() {
	var err error
	sc.newStdErr, err = os.CreateTemp("", "testErr")
	if err != nil {
		panic(err)
	}
	sc.oldStdErr = os.Stdout
	os.Stderr = sc.newStdErr
}

func (sc *StdChanger) Restore() {
	defer os.Remove(sc.newStdErr.Name())
	var err error
	sc.stdErrOutput, err = ioutil.ReadFile(sc.newStdErr.Name())
	if err != nil {
		panic(err)
	}
	os.Stderr = sc.oldStdErr
}

func (sc *StdChanger) Stderr() string {
	if sc.stdErrOutput == nil {
		panic(errors.New("You need to restore first"))
	}
	return string(sc.stdErrOutput)
}

func sendToAll(l Logger) {
	l.Warning(msgs[WarMsg])
	l.Info(msgs[InfoMsg])
	l.Debug(errors.New(msgs[DebMsg]))
	l.Error(errors.New(msgs[ErrMsg]))
}

func TestStdChanger(t *testing.T) {
	c := new(StdChanger)
	c.Redirect()
	fmt.Fprintln(os.Stderr, "hello")
	c.Restore()
	if !strings.Contains(c.Stderr(), "hello") {
		t.Log("Stderr:", c.Stderr())
		t.Error("Didn't find redirected message")
	}
}

func TestWithDebug(t *testing.T) {
	c := new(StdChanger)
	c.Redirect()
	l := NewLogger(Debug)
	if !l.doDebug {
		t.Error("Debug not enabled but should be")
	}
	sendToAll(l)
	c.Restore()
	for k, m := range msgs {
		if !strings.Contains(c.Stderr(), m) {
			t.Log("Stderr:", c.Stderr())
			t.Errorf("Stderr does not contain %v message: %v", k.Name(), m)
		}
	}
}

func TestNoFlags(t *testing.T) {
	c := new(StdChanger)
	c.Redirect()
	l := NewLogger(0)
	sendToAll(l)
	c.Restore()
	if l.doDebug {
		t.Error("Debug set when it shouldn't be")
	}
	for k, m := range msgs {
		if k == DebMsg {
			if strings.Contains(c.Stderr(), m) {
				t.Log("Stderr:", c.Stderr())
				t.Errorf("Stderr incorrectly contains %v message: %v", k.Name(), m)
			}
		} else if !strings.Contains(c.Stderr(), m) {
			t.Log("Stderr:", c.Stderr())
			t.Errorf("Stderr does not contain %v message: %v", k.Name(), m)
		}
	}
}

func TestEmpty(t *testing.T) {
	c := new(StdChanger)
	c.Redirect()
	l := NewLogger(Empty)
	sendToAll(l)
	c.Restore()
	if c.Stderr() != "" {
		t.Log("Stderr:", c.Stderr())
		t.Error("Empty logger outputted something even though it shouldn't")
	}
}

func TestQuiet(t *testing.T) {
	c := new(StdChanger)
	c.Redirect()
	l := NewLogger(Quiet)
	sendToAll(l)
	c.Restore()
	for k, m := range msgs {
		if k != ErrMsg {
			if strings.Contains(c.Stderr(), m) {
				t.Log("Stderr:", c.Stderr())
				t.Errorf("Stderr incorrectly contains %v message: %v", k.Name(), m)
			}
		} else if !strings.Contains(c.Stderr(), m) {
			t.Log("Stderr:", c.Stderr())
			t.Errorf("Stderr does not contain %v message: %v", k.Name(), m)
		}
	}
}
