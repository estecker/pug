package logs

import (
	"github.com/leg100/pug/internal/tui"
)

func coloredLogLevel(level string) string {
	levelColor := tui.InfoLogLevel
	switch level {
	case "ERROR":
		levelColor = tui.ErrorLogLevel
	case "WARN":
		levelColor = tui.WarnLogLevel
	case "DEBUG":
		levelColor = tui.DebugLogLevel
	case "INFO":
		levelColor = tui.InfoLogLevel
	}
	return tui.Bold.Foreground(levelColor).Render(level)
}
