package easylog

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

const (
	info  = "INFO"
	debug = "DEBUG"
	err   = "ERROR"
	fatal = "FATAL"
)

type Logger interface {
	Info(string)
	Infof(string, ...any)
	Debug(string)
	Debugf(string, ...any)
	Error(any)
	Errorf(string, ...any)
	Fail(any)
	Failf(string, ...any)
}

type defaultLogger struct {
	writer    io.Writer
	env       string
	hideDebug string
	onlyError string
}

type DefaultLoggerOptions struct {
	Writer         io.Writer
	Env            string
	HideDebugInEnv string
	OnlyErrorInEnv string
}

func NewDefaultLogger(opt DefaultLoggerOptions) Logger {
	return &defaultLogger{
		writer:    opt.Writer,
		env:       strings.ToUpper(opt.Env),
		hideDebug: strings.ToUpper(opt.HideDebugInEnv),
		onlyError: strings.ToUpper(opt.OnlyErrorInEnv),
	}
}

func (logger *defaultLogger) Info(msg string) {
	if logger.shouldHideNonErr() {
		return
	}
	logger.log(info, msg)
}

func (logger *defaultLogger) Infof(msg string, args ...any) {
	if logger.shouldHideNonErr() {
		return
	}
	logger.log(info, fmt.Sprintf(msg, args...))
}

func (logger *defaultLogger) Debug(msg string) {
	if logger.shouldHideNonErr() {
		return
	}
	if logger.shouldHideDebug() {
		return
	}
	logger.log(debug, msg)
}

func (logger *defaultLogger) Debugf(msg string, args ...any) {
	if logger.shouldHideNonErr() {
		return
	}
	if logger.shouldHideDebug() {
		return
	}
	logger.log(debug, fmt.Sprintf(msg, args...))
}

func (logger *defaultLogger) Error(msg any) {
	logger.log(err, fmt.Sprintf("%v", msg))
}

func (logger *defaultLogger) Errorf(msg string, args ...any) {
	logger.log(err, fmt.Sprintf(msg, args...))
}

func (logger *defaultLogger) Fail(msg any) {
	logger.log(fatal, fmt.Sprintf("%v", msg))
	os.Exit(1)
}

func (logger *defaultLogger) Failf(msg string, args ...any) {
	logger.log(fatal, fmt.Sprintf(msg, args...))
	os.Exit(1)
}

func (logger *defaultLogger) log(mode, msg string) {
	prefix := fmt.Sprintf("%s | %s | %s:", time.Now().Format(time.RFC3339), logger.env, mode)
	fmt.Fprintf(logger.writer, "%s %s\n", prefix, msg)
}

func (logger *defaultLogger) shouldHideNonErr() bool {
	return logger.env == logger.onlyError
}

func (logger *defaultLogger) shouldHideDebug() bool {
	return logger.env == logger.hideDebug
}
