package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

type Option struct {
	FileName         string
	Path             string
	Level            string
	MaxFileSize      int
	MaxTotalFile     int
	RetentionPeriode int
	IsAddSource      bool
}

func Init(opt Option) error {
	if opt.FileName == "" {
		opt.FileName = "app"
	}

	if opt.Path == "" {
		opt.Path = fmt.Sprintf("%s/logs", os.TempDir())
	}

	if opt.Level == "" {
		opt.Level = "info"
	}

	if opt.RetentionPeriode == 0 {
		opt.RetentionPeriode = 7
	}

	if opt.MaxFileSize == 0 {
		opt.MaxFileSize = 500
	}

	if opt.MaxTotalFile == 0 {
		opt.MaxTotalFile = 7
	}

	logRotator := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s.log", opt.Path, opt.FileName),
		MaxSize:    opt.MaxFileSize, // megabytes
		MaxBackups: opt.MaxTotalFile,
		MaxAge:     opt.RetentionPeriode, //days
	}

	logLvl, err := getLogLevel(opt.Level)
	if err != nil {
		return err
	}

	handler := slog.NewJSONHandler(
		io.MultiWriter(os.Stdout, logRotator),
		&slog.HandlerOptions{
			AddSource: opt.IsAddSource,
			Level:     logLvl,
		},
	)

	slog.SetDefault(slog.New(handler))
	return nil
}

func getLogLevel(logLevel string) (slog.Level, error) {
	lvl := strings.ToLower(logLevel)
	switch lvl {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return -10, fmt.Errorf("invalid log level: %s", lvl)
	}
}
