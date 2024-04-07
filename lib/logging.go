package lib

import (
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"os"
)

const loggingTemplate = "{{datetime}} {{level}} [{{caller}}] {{message}} {{data}} {{extra}}\n"

var logger = InitLogger()

func InitLogger() *slog.Logger {
	formatter := slog.NewTextFormatter()
	formatter.EnableColor = true
	formatter.ColorTheme = slog.ColorTheme
	formatter.TimeFormat = "15:04:05.000"
	formatter.SetTemplate(loggingTemplate)
	h := handler.NewConsoleHandler(slog.AllLevels)
	h.SetFormatter(formatter)
	logger := slog.NewWithHandlers(h)
	logger.ExitFunc = os.Exit
	return logger
}
