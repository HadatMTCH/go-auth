package log

import (
	"base-api/constants"
	"base-api/infra/log_rotator"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

var (
	rotator = &log_rotator.Logger{
		Filename: constants.LogFile,
	}
	Logger *logrus.Logger
)

func InitializeLogger(appInfo *constants.AppInfo) {
	Logger = logrus.New()

	// Set output to rotator and stdout
	Logger.SetOutput(rotator)
	Logger.AddHook(newStdoutHook())

	Logger.SetReportCaller(true)
	Logger.SetFormatter(&logFormat{
		TextFormatter: logrus.TextFormatter{
			DisableColors: true,
		},
	})

	// Set log level from environment
	setLogLevel()

	// Add application context
	Logger.WithFields(logrus.Fields{
		"service.name":    appInfo.AppName,
		"service.version": appInfo.AppVersion,
		"service.commit":  appInfo.AppCommit,
		"environment":     os.Getenv("ENV"),
	})
}

type logFormat struct {
	logrus.TextFormatter
}

func (f *logFormat) Format(entry *logrus.Entry) ([]byte, error) {
	var buffer strings.Builder

	// Base log structure
	buffer.WriteString(fmt.Sprintf("[%s]", entry.Time.Format(constants.LogDateFormatWithTime)))
	buffer.WriteString(fmt.Sprintf("[%s]", strings.ToUpper(entry.Level.String())))

	// Request context
	if traceID, ok := entry.Data["trace_id"]; ok {
		buffer.WriteString(fmt.Sprintf("[TraceID:%v]", traceID))
	}

	// HTTP context
	if method, ok := entry.Data["http.method"]; ok {
		buffer.WriteString(fmt.Sprintf("[%s]", method))
	}
	if path, ok := entry.Data["http.path"]; ok {
		buffer.WriteString(fmt.Sprintf("[%s]", path))
	}
	if status, ok := entry.Data["http.status"]; ok {
		buffer.WriteString(fmt.Sprintf("[Status:%v]", status))
	}

	// Message
	buffer.WriteString(fmt.Sprintf(" %s", entry.Message))

	// Additional fields
	if entry.Data["error"] != nil {
		buffer.WriteString(fmt.Sprintf("\nERROR: %v", entry.Data["error"]))
	}
	if entry.Data["duration"] != nil {
		buffer.WriteString(fmt.Sprintf("\nDURATION: %v", entry.Data["duration"]))
	}

	// Caller information
	if entry.HasCaller() {
		buffer.WriteString(fmt.Sprintf("\nFILE: %s:%d",
			entry.Caller.File,
			entry.Caller.Line,
		))
	}

	buffer.WriteString("\n\n")
	return []byte(buffer.String()), nil
}

func setLogLevel() {
	switch os.Getenv("LOG_LEVEL") {
	case "trace":
		Logger.SetLevel(logrus.TraceLevel)
	case "debug":
		Logger.SetLevel(logrus.DebugLevel)
	case "warn":
		Logger.SetLevel(logrus.WarnLevel)
	case "error":
		Logger.SetLevel(logrus.ErrorLevel)
	default:
		Logger.SetLevel(logrus.InfoLevel)
	}
}

// Helper to create request-aware logger
func FromContext(c echo.Context) *logrus.Entry {
	fields := logrus.Fields{}

	if c != nil {
		fields = logrus.Fields{
			"trace_id":    c.Get("trace_id"),
			"http.method": c.Request().Method,
			"http.path":   c.Path(),
			"http.status": c.Response().Status,
			"client_ip":   c.RealIP(),
		}
	}

	return Logger.WithFields(fields)
}

// Stdout hook for dual logging
type stdoutHook struct{}

func newStdoutHook() *stdoutHook             { return &stdoutHook{} }
func (h *stdoutHook) Levels() []logrus.Level { return logrus.AllLevels }
func (h *stdoutHook) Fire(entry *logrus.Entry) error {
	line, _ := entry.String()
	os.Stdout.Write([]byte(line))
	return nil
}
