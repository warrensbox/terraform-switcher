package lib

import (
	"os"

	"github.com/gookit/color"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/mattn/go-isatty"
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

func isColorLogging() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	} else if color.SupportColor() {
		if os.Getenv("FORCE_COLOR") == "" {
			return isatty.IsTerminal(os.Stdout.Fd())
		}
		return true
	}
	return false
}

func NewStderrConsoleWithLF(lf slog.LevelFormattable) *handler.ConsoleHandler {
	h := handler.NewIOWriterWithLF(os.Stderr, lf)

	// default use text formatter
	f := slog.NewTextFormatter()
	// default enable color on console
	f.WithEnableColor(isColorLogging())
	h.SetFormatter(f)
	return h
}

func NewStderrConsoleHandler(levels []slog.Level) *handler.ConsoleHandler {
	return NewStderrConsoleWithLF(slog.NewLvsFormatter(levels))
}

func InitLogger(logLevel string) *slog.Logger {
	formatter := slog.NewTextFormatter()
	formatter.EnableColor = isColorLogging()
	formatter.ColorTheme = slog.ColorTheme
	formatter.TimeFormat = "15:04:05.000"

	var h *handler.ConsoleHandler
	switch logLevel {
	case "ERROR":
		h = NewStderrConsoleHandler(ErrorLogging)
		formatter.SetTemplate(loggingTemplateDebug)
	case "TRACE":
		h = NewStderrConsoleHandler(TraceLogging)
		formatter.SetTemplate(loggingTemplateDebug)
	case "DEBUG":
		h = NewStderrConsoleHandler(DebugLogging)
		formatter.SetTemplate(loggingTemplateDebug)
	case "NOTICE":
		h = NewStderrConsoleHandler(NoticeLogging)
		formatter.SetTemplate(loggingTemplate)
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
