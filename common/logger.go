package common

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type Config struct {
	Filepath string
}

type logger struct {
	*logrus.Logger
}

var config *Config
var log *logger

func InitLoggerConfig(conf *Config) error {
	if config != nil {
		return fmt.Errorf("logger config already initialized")
	}
	config = conf
	return nil
}

func (l *logger) init() error {
	l.Logger = logrus.New()
	l.Level = logrus.WarnLevel

	if config == nil {
		return fmt.Errorf("config is not set")
	}

	dirPath := filepath.Dir(config.Filepath)

	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory for log path: %v", err)
	}

	info, err := os.Stat(dirPath)
	if err != nil {
		return fmt.Errorf("failed to stat directory: %v", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("the path is not a directory: %v", dirPath)
	}

	file, err := os.OpenFile(config.Filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	l.Out = file

	l.Formatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyFile:  "file",
		},
	}
	return nil
}

func GetLogger() (*logrus.Logger, error) {
	if log == nil {
		log = new(logger)
		err := log.init()
		if err != nil {
			return nil, err
		}
	}
	return log.Logger, nil

}
