package logger

import (
	"flag"
	"net/http"
	"os"
	"regexp"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	_atom   = zap.NewAtomicLevel()
	_logger *zap.SugaredLogger
)

func init() {
	var (
		lvl            = zap.InfoLevel
		levelFlagName  = "log-level"
		levelFlagUsage = "minimum enabled logging level. debug|info|warn|error|dpanic|panic|fatal"
	)

	flag.Var(&lvl, levelFlagName, levelFlagUsage)

	logArgs := regexp.MustCompile(`-{1,2}` + levelFlagName + `(?:\s+|\s*=\s*)(\w+)`).
		FindString(strings.Join(os.Args[1:], " "))

	if logArgs != "" {
		// use local FlagSet to parse immediately
		flagSet := flag.NewFlagSet("logger", flag.ContinueOnError)
		flagSet.Var(&lvl, levelFlagName, levelFlagUsage)
		flagSet.Parse(regexp.MustCompile(`\s+`).Split(logArgs, 2))
		_atom.SetLevel(lvl)
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = ""

	_logger = zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		_atom,
	)).Sugar()
}

func Debug(msg string, kvs ...interface{}) {
	_logger.Debugw(msg, kvs...)
}

func Debugf(msg string, args ...interface{}) {
	_logger.Debugf(msg, args...)
}

func Info(msg string, kvs ...interface{}) {
	_logger.Infow(msg, kvs...)
}

func Infof(msg string, args ...interface{}) {
	_logger.Infof(msg, args...)
}

func Warn(msg string, kvs ...interface{}) {
	_logger.Warnw(msg, kvs...)
}

func Warnf(msg string, args ...interface{}) {
	_logger.Warnf(msg, args...)
}

func Error(msg string, kvs ...interface{}) {
	_logger.Errorw(msg, kvs...)
}

func Errorf(msg string, args ...interface{}) {
	_logger.Errorf(msg, args...)
}

func DPanic(msg string, kvs ...interface{}) {
	_logger.DPanicw(msg, kvs...)
}

func DPanicf(msg string, args ...interface{}) {
	_logger.DPanicf(msg, args...)
}

func Panic(msg string, kvs ...interface{}) {
	_logger.Panicw(msg, kvs...)
}

func Panicf(msg string, args ...interface{}) {
	_logger.Panicf(msg, args...)
}

func Fatal(msg string, kvs ...interface{}) {
	_logger.Fatalw(msg, kvs...)
}

func Fatalf(msg string, args ...interface{}) {
	_logger.Fatalf(msg, args...)
}

func HttpHandler() http.Handler {
	return _atom
}

func Sync() {
	_logger.Sync()
}
