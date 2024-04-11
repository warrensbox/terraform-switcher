package lib

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"os"
)

const loggingTemplate = "{{datetime}} {{level}} [{{caller}}] {{message}} {{data}} {{extra}}\n"

var logger *slog.Logger

func InitLogger(logLevel string) *slog.Logger {
	formatter := slog.NewTextFormatter()
	formatter.EnableColor = true
	formatter.ColorTheme = slog.ColorTheme
	formatter.TimeFormat = "15:04:05.000"
	formatter.SetTemplate(loggingTemplate)

	var h *handler.ConsoleHandler
	if logLevel == "TRACE" {
		h = handler.NewConsoleHandler(TraceLogging)
	} else if logLevel == "DEBUG" {
		h = handler.NewConsoleHandler(DebugLogging)
	} else if logLevel == "NOTICE" {
		h = handler.NewConsoleHandler(NoticeLogging)
	} else {
		h = handler.NewConsoleHandler(NormalLogging)
	}

	h.SetFormatter(formatter)
	newLogger := slog.NewWithHandlers(h)
	newLogger.ExitFunc = os.Exit
	logger = newLogger
	return newLogger
}

var (
	NormalLogging = slog.Levels{slog.PanicLevel, slog.FatalLevel, slog.ErrorLevel, slog.WarnLevel, slog.InfoLevel}
	NoticeLogging = slog.Levels{slog.PanicLevel, slog.FatalLevel, slog.ErrorLevel, slog.WarnLevel, slog.InfoLevel, slog.NoticeLevel}
	DebugLogging  = slog.Levels{slog.PanicLevel, slog.FatalLevel, slog.ErrorLevel, slog.WarnLevel, slog.InfoLevel, slog.NoticeLevel, slog.DebugLevel}
	TraceLogging  = slog.Levels{slog.PanicLevel, slog.FatalLevel, slog.ErrorLevel, slog.WarnLevel, slog.InfoLevel, slog.NoticeLevel, slog.DebugLevel, slog.TraceLevel}
)
