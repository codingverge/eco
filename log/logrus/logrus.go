package logrus

import "github.com/sirupsen/logrus"

type (
	Logrus struct {
		*logrus.Entry
		version string
	}
	options struct {
		l             *logrus.Logger
		level         *logrus.Level
		formatter     logrus.Formatter
		format        string
		reportCaller  bool
		exitFunc      func(int)
		leakSensitive bool
		redactionText string
		hooks         []logrus.Hook
	}
	Option func(*options)
)
