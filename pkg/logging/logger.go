package logging

import (
	"flag"
	"github.com/sirupsen/logrus"
	"os"
)

type Logger struct {
	*logrus.Entry
}

func New() *Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{})
	l.SetOutput(os.Stdout)

	l.SetLevel(logrus.InfoLevel)
	var debug bool
	flag.BoolVar(&debug, "debug", false, "enable debug mode")
	flag.Parse()
	if debug == true {
		l.SetLevel(logrus.DebugLevel)
	}

	return &Logger{logrus.NewEntry(l)}
}
