package logger

import (
	"encoding/json"
	"fmt"
	"gopa/config"
	"os"
	"path/filepath"
	"time"

	"github.com/YueHonghui/rfw"
	"github.com/sirupsen/logrus"
)

var (
	Log           *logrus.Logger
	AccessLog     *logrus.Logger
	RuntimeLog    *logrus.Logger
	RuntimeErrLog *logrus.Logger
	MetricsLog    *logrus.Logger
)

const (
	LogRemain int = 10
	LogDir        = "/var/log/gopa"
)

type logFieldKey string

type MetricsJSONFormatter struct{}

//Format log fotmat
func (f *MetricsJSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Note this doesn't include Time, Level and Message which are available on
	// the Entry.
	entry.Data["time"] = entry.Time.Format(time.RFC3339)
	serialized, err := json.Marshal(entry.Data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}

//MetricsEmit logs the metrics message
func MetricsEmit(method, reqID string, message interface{}, success bool) {
	MetricsLog.WithFields(logrus.Fields{
		"topic":   "trace",
		"method":  method,
		"reqID":   reqID,
		"success": success,
	}).Info(message)
}

//RuntimeEmit logs the runtime message
func RuntimeEmit(method, reqID string, message interface{}, success bool) {
	RuntimeErrLog.WithFields(logrus.Fields{
		"topic":   "trace",
		"method":  method,
		"reqID":   reqID,
		"success": success,
	}).Warn(message)
}

//SetLog init the logger config
func SetLog() error {
	logrus.SetFormatter(&logrus.JSONFormatter{}) // Log as JSON instead of the default ASCII formatter.
	logrus.SetOutput(os.Stdout)                  // Output to stdout instead of the default stderr, Can be any io.Writer
	logrus.SetLevel(logrus.TraceLevel)           // Only log the warning severity or above.

	logDir := config.GetConfig().Log.Dir
	logRemain := config.GetConfig().Log.Remain
	// Logging with rwf
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		logrus.WithError(err).Fatalf("Failed mkdir -p %s", logDir)
	}

	if rfw, err := rfw.NewWithOptions(filepath.Join(logDir, "gopa-access"), rfw.WithCleanUp(logRemain)); err != nil {
		AccessLog = logrus.StandardLogger()
	} else {
		AccessLog = &logrus.Logger{
			Out:       rfw,
			Level:     logrus.InfoLevel,
			Formatter: &logrus.JSONFormatter{},
		}
	}

	if rfw, err := rfw.NewWithOptions(filepath.Join(logDir, "gopa-runtime"), rfw.WithCleanUp(logRemain)); err != nil {
		RuntimeErrLog = logrus.StandardLogger()
	} else {
		RuntimeErrLog = &logrus.Logger{
			Out:       rfw,
			Level:     logrus.DebugLevel,
			Formatter: &logrus.JSONFormatter{},
		}
	}

	if rfw, err := rfw.NewWithOptions(filepath.Join(logDir, "gopa-metrics"), rfw.WithCleanUp(logRemain)); err != nil {
		MetricsLog = logrus.StandardLogger()
	} else {
		MetricsLog = &logrus.Logger{
			Out:       rfw,
			Level:     logrus.InfoLevel,
			Formatter: &logrus.JSONFormatter{},
		}
	}
	return nil
}
