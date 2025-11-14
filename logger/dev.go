package logger

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type DevLogger struct {
	l *log.Logger
}

func NewDevLogger() *DevLogger {
	l := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    false,
		TimeFormat:      "15:04:05",
		ReportTimestamp: true,
	})

	styles := log.DefaultStyles()

	styles.Levels[log.DebugLevel] = lipgloss.NewStyle().
		SetString(strings.ToUpper(log.DebugLevel.String())).
		Bold(true).
		Foreground(lipgloss.Color("#5FD1F9"))

	styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
		SetString(strings.ToUpper(log.InfoLevel.String())).
		Bold(true).
		Foreground(lipgloss.Color("#00D75F"))

	styles.Levels[log.WarnLevel] = lipgloss.NewStyle().
		SetString(strings.ToUpper(log.WarnLevel.String())).
		Bold(true).
		Foreground(lipgloss.Color("#FFD75F"))

	styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
		SetString(strings.ToUpper(log.ErrorLevel.String())).
		Bold(true).
		Foreground(lipgloss.Color("#FF5F5F"))

	styles.Timestamp = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	styles.Prefix = lipgloss.NewStyle().Foreground(lipgloss.Color("#6666FF"))
	styles.Key = lipgloss.NewStyle().Foreground(lipgloss.Color("#A8A8A8"))
	styles.Value = lipgloss.NewStyle().Foreground(lipgloss.Color("#EAEAEA"))

	l.SetStyles(styles)

	return &DevLogger{l: l}
}

func (d *DevLogger) Info(msg string, kv ...any)  { d.l.Info(msg, kv...) }
func (d *DevLogger) Warn(msg string, kv ...any)  { d.l.Warn(msg, kv...) }
func (d *DevLogger) Error(msg string, kv ...any) { d.l.Error(msg, kv...) }
func (d *DevLogger) Debug(msg string, kv ...any) { d.l.Debug(msg, kv...) }

func (d *DevLogger) Infof(template string, args ...any)  { d.l.Infof(template, args...) }
func (d *DevLogger) Warnf(template string, args ...any)  { d.l.Warnf(template, args...) }
func (d *DevLogger) Errorf(template string, args ...any) { d.l.Errorf(template, args...) }
func (d *DevLogger) Debugf(template string, args ...any) { d.l.Debugf(template, args...) }
