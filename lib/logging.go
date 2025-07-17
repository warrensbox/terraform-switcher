package lib

import (
	"os"
	"slices"
	"strings"

	"github.com/gookit/color"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/mattn/go-isatty"
)

var (
	logger               *slog.Logger
	loggingTemplateDebug = "{{datetime}} {{level}} [{{caller}}] {{message}} {{data}} {{extra}}\n"
	loggingTemplate      = "{{datetime}} {{level}} {{message}} {{data}} {{extra}}\n"
	// Parent lib: https://github.com/gookit/slog/blob/f857defd050dd7fc3c3013134cf50ed51b917a1f/common.go#L69-L88
	loggingLevel = map[string]slog.Levels{
		// Special case to disable (suppress) logging
		"OFF": {},
		// High severity, unrecoverable errors (internally calls `panic()`)
		"PANIC": {slog.PanicLevel},
		// Fatal, unrecoverable errors
		"FATAL": {slog.PanicLevel, slog.FatalLevel},
		// Runtime errors that should definitely be noted
		"ERROR": {slog.PanicLevel, slog.FatalLevel, slog.ErrorLevel},
		// Non-critical entries that deserve eyes
		"WARN": {slog.PanicLevel, slog.FatalLevel, slog.ErrorLevel, slog.WarnLevel},
		// Default log level, messages that highlight the progress
		"INFO": {slog.PanicLevel, slog.FatalLevel, slog.ErrorLevel, slog.WarnLevel, slog.InfoLevel},
		// Normal operational entries, but not necessarily noteworthy
		"NOTICE": {slog.PanicLevel, slog.FatalLevel, slog.ErrorLevel, slog.WarnLevel, slog.InfoLevel, slog.NoticeLevel},
		// Verbose logging, useful for developers
		"DEBUG": {slog.PanicLevel, slog.FatalLevel, slog.ErrorLevel, slog.WarnLevel, slog.InfoLevel, slog.NoticeLevel, slog.DebugLevel},
		// Even more finer-grained informational events
		"TRACE": {slog.PanicLevel, slog.FatalLevel, slog.ErrorLevel, slog.WarnLevel, slog.InfoLevel, slog.NoticeLevel, slog.DebugLevel, slog.TraceLevel},
	}
	logUnknownLogLevel bool
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

func LogLevels() []string {
	levels := make([]string, 0, len(loggingLevel))
	for level := range loggingLevel {
		levels = append(levels, level)
	}
	slices.Sort(levels) // Sort the result as we also use it in output messages
	return levels
}

func InitLogger(logLevel string) *slog.Logger {
	// Allow lower or mixed case for log levels to
	// provide flexibility and improve user experience
	originalLogLevelValue := logLevel
	fallbackLogLevel := "INFO"
	logLevel = strings.ToUpper(logLevel)

	formatter := slog.NewTextFormatter()
	formatter.ColorTheme = slog.ColorTheme
	formatter.EnableColor = isColorLogging()
	formatter.SetTemplate(loggingTemplate)
	formatter.TimeFormat = "15:04:05.000"

	useDebugTemplateLogLevels := []string{"DEBUG", "ERROR", "FATAL", "PANIC", "TRACE"}

	var h *handler.ConsoleHandler
	// Safe log level fallback just in case
	// See `initParams()` in lib/param_parsing/parameters.go for default log level
	h = NewStderrConsoleHandler(loggingLevel[fallbackLogLevel])

	isUnknownLogLevel := false

	if slices.Contains(LogLevels(), logLevel) {
		h = NewStderrConsoleHandler(loggingLevel[logLevel])
		if slices.Contains(useDebugTemplateLogLevels, logLevel) {
			formatter.SetTemplate(loggingTemplateDebug)
		}
	} else {
		isUnknownLogLevel = true
	}

	h.SetFormatter(formatter)
	newLogger := slog.NewWithHandlers(h)
	newLogger.ExitFunc = os.Exit
	logger = newLogger

	if isUnknownLogLevel {
		// TODO: [20250717] Drop the below `if`-conditional and switch to `logger.Fatalf()` to fail
		// on unknown log level, say, in a couple of months so that users have enough time to notice
		// this warning and fix their configuration (if any)
		// logger.Fatalf("Unhandled logging level: %q (must be one of: %s)", originalLogLevelValue, strings.Join(LogLevels(), ", "))
		if !logUnknownLogLevel { // Only log this warning once on very first occurrence of unknown log level
			logUnknownLogLevel = true
			logger.Warnf(
				"Unhandled logging level: %q (must be one of: %s). Falling back to %s\n\t!!! THIS WILL BE A FATAL ERROR IN THE FUTURE !!!",
				originalLogLevelValue,
				strings.Join(LogLevels(), ", "),
				fallbackLogLevel,
			)
		}
	}

	return logger
}
