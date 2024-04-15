package tlog

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	slogmulti "github.com/samber/slog-multi"
)

type Tlog struct {
	*slog.Logger
}

func GetLogger(level slog.Level) *Tlog {
	projDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	logDir := filepath.Join(projDir, "logs")

	logFilename := filepath.Join(logDir, time.Now().Format("2006-01-02")+".log")

	logFile, err := os.OpenFile(logFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Errorf("Failed to open file '%s' %w", logFilename, err))
	}

	opts := &slog.HandlerOptions{
		AddSource:   true,
		Level:       level,
		ReplaceAttr: nil,
	}
	fileHandler := slog.NewJSONHandler(logFile, opts)

	outHandler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(slogmulti.Fanout(fileHandler, outHandler))
	slog.SetDefault(logger)

	return &Tlog{logger}
}
