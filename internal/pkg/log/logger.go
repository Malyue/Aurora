package log

import (
	_const "Aurora/internal/pkg/const"
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
)

type Config struct {
	LogPath     string `yaml:"logPath"`
	Mode        string `yaml:"mode"`
	ServiceName string `yaml:"serviceName"`
}

const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

type LogFormatter struct {
}

type DefaultFieldHook struct {
	ServiceName string
}

func (hook *DefaultFieldHook) Fire(entry *logrus.Entry) error {
	entry.Data["Service"] = hook.ServiceName
	return nil
}

func (hook *DefaultFieldHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func InitLogger(cfg *Config) *logrus.Logger {
	logger := logrus.New()
	file, _ := os.OpenFile(cfg.LogPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if cfg.Mode == _const.DEBUG || cfg.Mode == _const.DEV {
		logger.SetOutput(io.MultiWriter(os.Stdout))
		logger.SetLevel(logrus.DebugLevel)
	} else if cfg.Mode == _const.PRODUCTION {
		logger.SetOutput(io.MultiWriter(os.Stdout, file))
		logger.SetLevel(logrus.ErrorLevel)
		logrus.SetOutput(io.MultiWriter(os.Stdout, file))
	}
	logger.SetReportCaller(true)
	logger.SetFormatter(&LogFormatter{})
	hook := &DefaultFieldHook{
		ServiceName: cfg.ServiceName,
	}
	logger.AddHook(hook)

	return logger
}

func (f LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format(_const.TIMESTAMPFROMAT)
	if entry.HasCaller() {
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)
		_, _ = fmt.Fprintf(b, "[%s] \033[%dm[%s]\033[0m %s %s %s \n]", timestamp, levelColor, entry.Level, fileVal, funcVal, entry.Message)
	} else {
		_, _ = fmt.Fprintf(b, "[%s] \033[%dm[%s]\033[0m %s %s %s \n]", timestamp, levelColor, entry.Level, entry.Message)
	}
	return b.Bytes(), nil
}
