package easylog

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

type writerMock struct {
	val []byte
}

func (w *writerMock) Write(b []byte) (int, error) {
	w.val = b
	return len(b), nil
}

func (w *writerMock) get() string {
	return string(w.val)
}

var writer writerMock

func getLogger(env string, hideDbg, showErr bool) Logger {
	dbg := ""
	if hideDbg {
		dbg = env
	}

	err := ""
	if showErr {
		err = env
	}

	return NewDefaultLogger(DefaultLoggerOptions{
		Env:            env,
		Writer:         &writer,
		HideDebugInEnv: dbg,
		OnlyErrorInEnv: err,
	})
}

func Test_defaultLogger(t *testing.T) {
	log := getLogger("dev", false, false)
	log.Info("test")

	if writer.get() == "" {
		t.Error("no bytes have been written by the logger")
	}

	log.Infof("test %d", 2)
	got := stripMessage(writer.get(), 2)
	if got != "INFO: test 2" {
		t.Errorf("incorrect message format: %s", got)
	}

	log.Debug("test")
	got = stripMessage(writer.get(), 2)
	if got != "DEBUG: test" {
		t.Errorf("incorrect message format: %s", got)
	}
	log.Debugf("test %d", 2)

	log2 := getLogger("PROD", true, false)
	log2.Debugf("skip this %s", "message")
	log2.Debug("skip this also")
	got = stripMessage(writer.get(), 2)
	if got != "DEBUG: test 2" {
		t.Errorf("incorrect message format: %s", got)
	}

	log2.Error(fmt.Errorf("passing an error"))
	got = stripMessage(writer.get(), 1)
	if got != "PROD" {
		t.Errorf("incorrect message format: %s", got)
	}

	log3 := getLogger("STAGE", false, true)
	log3.Errorf("error %d", 2)
	log3.Info("skip this")
	log3.Infof("skip this %d", 2)
	log3.Debug("skip this")
	log3.Debugf("skip this %d", 2)
	got = stripMessage(writer.get(), 2)
	if got != "ERROR: error 2" {
		t.Errorf("incorrect message format: %s", got)
	}
}

func TestCrash(t *testing.T) {
	log := getLogger("dev", false, false)
	if os.Getenv("BE_CRASHER") == "1" {
		log.Fail("fail")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestCrash")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestCrash2(t *testing.T) {
	log := getLogger("dev", false, false)
	if os.Getenv("BE_CRASHER") == "2" {
		log.Failf("fail %s", "again")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestCrash2")
	cmd.Env = append(os.Environ(), "BE_CRASHER=2")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func stripMessage(msg string, i int) string {
	s := strings.Split(msg, "|")
	if len(s) < (i + 1) {
		return ""
	}

	return strings.TrimSpace(s[i])
}
