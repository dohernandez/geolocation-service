package log

import (
	logging "github.com/hellofresh/logging-go"
	"github.com/sirupsen/logrus"
)

// NewLog creates Logger
func NewLog(config logging.LogConfig) (*logrus.Logger, error) {
	if config.Level == "" {
		config.Level = "debug"
	}
	err := config.Apply()
	if err != nil {
		return nil, err
	}

	l := logrus.New()
	l.Formatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "timestamp",
			logrus.FieldKeyMsg:  "message",
		},
	}

	ll := logrus.GetLevel()
	l.SetLevel(ll)

	return l, nil
}
