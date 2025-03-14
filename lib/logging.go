package lib

import (
	"os"

	"github.com/gookit/color"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

var (
	loggingTemplateDebug = "{{datetime}} {{level}} [{{caller}}] {{message}} {{data}} {{extra}}\n"
	loggingTemplate      = "{{datetime}} {{level}} {{message}} {{data}} {{extra}}\n"
	logger               *slog.Logger
	ErrorLogging         = slog.Levels{slog.PanicLevel, slog.FatalLevel, slog.ErrorLevel}
	NormalLogging        = slog.Levels{slog.PanicLevel, slog.FatalLevel, slog.ErrorLevel, slog.WarnLevel, slog.InfoLevel}
	NoticeLogging        = slog.Levels{slog.PanicLevel, slog.FatalLevel, slog.ErrorLevel, slog.WarnLevel, slog.InfoLevel, slog.NoticeLevel}
	DebugLogging         = slog.Levels{slog.PanicLevel, slog.FatalLevel, slog.ErrorLevel, slog.WarnLevel, slog.InfoLevel, slog.NoticeLevel, slog.DebugLevel}
	TraceLogging         = slog.Levels{slog.PanicLevel, slog.FatalLevel, slog.ErrorLevel, slog.WarnLevel, slog.InfoLevel, slog.NoticeLevel, slog.DebugLevel, slog.TraceLevel}
)

func NewStderrConsoleWithLF(lf slog.LevelFormattable) *handler.ConsoleHandler {
	h := handler.NewIOWriterWithLF(os.Stderr, lf)

	// default use text formatter
	f := slog.NewTextFormatter()
	// default enable color on console
	f.WithEnableColor(color.SupportColor())

	h.SetFormatter(f)
	return h
}

func NewStderrConsoleHandler(levels []slog.Level) *handler.ConsoleHandler {
	return NewStderrConsoleWithLF(slog.NewLvsFormatter(levels))
}

func InitLogger(logLevel string) *slog.Logger {
	formatter := slog.NewTextFormatter()
	formatter.EnableColor = true
	formatter.ColorTheme = slog.ColorTheme
	formatter.TimeFormat = "15:04:05.000"

	var h *handler.ConsoleHandler
	switch logLevel {
	case "ERROR":
		h = NewStderrConsoleHandler(ErrorLogging)
		formatter.SetTemplate(loggingTemplateDebug)
		break
	case "TRACE":
		h = NewStderrConsoleHandler(TraceLogging)
		formatter.SetTemplate(loggingTemplateDebug)
		break
	case "DEBUG":
		h = NewStderrConsoleHandler(DebugLogging)
		formatter.SetTemplate(loggingTemplateDebug)
		break
	case "NOTICE":
		h = NewStderrConsoleHandler(NoticeLogging)
		formatter.SetTemplate(loggingTemplate)
		break
	default:
		h = NewStderrConsoleHandler(NormalLogging)
		formatter.SetTemplate(loggingTemplate)
	}

	h.SetFormatter(formatter)
	newLogger := slog.NewWithHandlers(h)
	newLogger.ExitFunc = os.Exit
	logger = newLogger
	return newLogger
}
