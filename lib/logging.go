package lib

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"os"
)

var (
	loggingTemplateDebug = "{{datetime}} {{level}} [{{caller}}] {{message}} {{data}} {{extra}}\n"
	loggingTemplate      = "{{datetime}} {{level}} {{message}} {{data}} {{extra}}\n"
	logger               *slog.Logger
	NormalLogging        = slog.Levels{slog.PanicLevel, slog.FatalLevel, slog.ErrorLevel, slog.WarnLevel, slog.InfoLevel}
	NoticeLogging        = slog.Levels{slog.PanicLevel, slog.FatalLevel, slog.ErrorLevel, slog.WarnLevel, slog.InfoLevel, slog.NoticeLevel}
	DebugLogging         = slog.Levels{slog.PanicLevel, slog.FatalLevel, slog.ErrorLevel, slog.WarnLevel, slog.InfoLevel, slog.NoticeLevel, slog.DebugLevel}
	TraceLogging         = slog.Levels{slog.PanicLevel, slog.FatalLevel, slog.ErrorLevel, slog.WarnLevel, slog.InfoLevel, slog.NoticeLevel, slog.DebugLevel, slog.TraceLevel}
)

func InitLogger(logLevel string) *slog.Logger {
	formatter := slog.NewTextFormatter()
	formatter.EnableColor = true
	formatter.ColorTheme = slog.ColorTheme
	formatter.TimeFormat = "15:04:05.000"

	var h *handler.ConsoleHandler
	switch logLevel {
	case "TRACE":
		h = handler.NewConsoleHandler(TraceLogging)
		formatter.SetTemplate(loggingTemplateDebug)
		break
	case "DEBUG":
		h = handler.NewConsoleHandler(DebugLogging)
		formatter.SetTemplate(loggingTemplateDebug)
		break
	case "NOTICE":
		h = handler.NewConsoleHandler(NoticeLogging)
		formatter.SetTemplate(loggingTemplateDebug)
		break
	default:
		h = handler.NewConsoleHandler(NormalLogging)
		formatter.SetTemplate(loggingTemplate)
	}

	h.SetFormatter(formatter)
	newLogger := slog.NewWithHandlers(h)
	newLogger.ExitFunc = os.Exit
	logger = newLogger
	return newLogger
}
