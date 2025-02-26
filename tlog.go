package tlog

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	slogmulti "github.com/samber/slog-multi"
)

func GetLogger(isDebug bool) *slog.Logger {
	projDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	logDir := filepath.Join(projDir, "logs")

	if _, err := os.Stat(logDir); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			panic(err)
		}
	}

	logFilename := filepath.Join(logDir, time.Now().Format("2006-01-02")+".log")

	logFile, err := os.OpenFile(logFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Errorf("Failed to open file '%s' %w", logFilename, err))
	}

	level := slog.LevelInfo
	if isDebug {
		level = slog.LevelDebug
	}

	replace := func(groups []string, a slog.Attr) slog.Attr {
		// Remove the directory from the source's filename.
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			source.File = filepath.Base(source.File)
		}
		return a
	}

	opts := &slog.HandlerOptions{
		AddSource:   true,
		Level:       level,
		ReplaceAttr: replace,
	}
	fileHandler := slog.NewJSONHandler(logFile, opts)

	outHandler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(slogmulti.Fanout(fileHandler, outHandler))
	slog.SetDefault(logger)

	return logger
}
