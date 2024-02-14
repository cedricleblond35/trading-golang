package config

import (
	"os"
	"regexp"

	"github.com/mdouchement/logger"
	"github.com/sirupsen/logrus"
)

// NewLogger returns a new logger.
// It panics if there is some troubles.
func NewLogger(name string) (log logger.Logger) {
	var err error
	l := logrus.New()

	l.Formatter = &logger.LogrusTextFormatter{
		PrefixRE:        regexp.MustCompile(`^(\[.*?\])\s`),
		ForceColors:     os.Getenv("LOG_FORCE_COLORS") == "enabled",
		ForceFormatting: os.Getenv("LOG_FORCE_FORMATTING") == "enabled",
		FullTimestamp:   os.Getenv("LOG_FORCE_FORMATTING") == "enabled",
	}

	l.Level, err = logrus.ParseLevel("fatal")
	if err != nil {
		panic(err)
	}

	// Global fields inside the Logrus mapper
	log = logger.WrapLogrus(l).WithField("logger_name", name)
	log.Infof("Log level: %s", l.Level)

	return
}
