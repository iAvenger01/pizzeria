package logging

import (
	"flag"
	"github.com/sirupsen/logrus"
	"os"
)

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func Init() {
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{})
	l.SetOutput(os.Stdout)

	l.SetLevel(logrus.InfoLevel)
	var debug bool
	flag.BoolVar(&debug, "debug", false, "enable debug mode")
	flag.Parse()
	if debug == true {
		l.SetLevel(logrus.DebugLevel)
	}

	e = logrus.NewEntry(l)
}

func GetLogger() Logger {
	return Logger{e}
}
