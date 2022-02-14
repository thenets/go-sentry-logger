package gosentry

import (
	"fmt"
	"strings"
	"time"

	sentry "github.com/getsentry/sentry-go"
	logrus "github.com/sirupsen/logrus"
)

type GoSentry struct {
}

func (s *GoSentry) Init(sentry_dsn string) error {
	logrus.Debug("Initializing Sentry...")
	logrus.Debug("Sentry DSN: " + sentry_dsn)

	// Start Sentry client
	err := sentry.Init(sentry.ClientOptions{
		Dsn: sentry_dsn,
	})
	if err != nil {
		logrus.Fatalf("sentry.Init: %s", err)
		return err
	} else {
		sentry_server_domain := strings.Split(strings.Split(sentry_dsn, "@")[1], "/")[0]
		logrus.Info("Sentry initialized: " + sentry_server_domain)
	}

	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)

	return nil
}

func (s *GoSentry) GeneralLogger(in interface{}, level logrus.Level, force_push_to_sentry bool) {
	var e error
	var msg string
	push_to_sentry := false

	// Extract the message from the interface
	switch v := in.(type) {
	case string:
		msg = v
	case error:
		msg = v.Error()
	default:
		fmt.Errorf("%v", e).Error()
	}

	// Switch on the log level
	switch level {
	case logrus.PanicLevel:
		defer logrus.Panic(msg)
		push_to_sentry = true
	case logrus.FatalLevel:
		defer logrus.Fatal(msg)
		push_to_sentry = true
	case logrus.ErrorLevel:
		sentry.CaptureException(e)
		logrus.Error(msg)
		push_to_sentry = true
	case logrus.WarnLevel:
		logrus.Warn(msg)
	case logrus.InfoLevel:
		logrus.Info(msg)
	case logrus.DebugLevel:
		logrus.Debug(msg)
	}

	// Push the message to Sentry
	if push_to_sentry || force_push_to_sentry {
		var sentry_event_id *sentry.EventID
		if e != nil {
			sentry_event_id = sentry.CaptureException(e)
		}
		if msg != "" {
			sentry_event_id = sentry.CaptureMessage(msg)
		}
		logrus.Debug("Sentry event ID: ", sentry_event_id)
		defer sentry.Flush(2 * time.Second)
	}
}

func (s *GoSentry) Panic(in interface{}) {
	s.GeneralLogger(in, logrus.PanicLevel, false)
}
func (s *GoSentry) Fatal(in interface{}) {
	s.GeneralLogger(in, logrus.FatalLevel, false)
}
func (s *GoSentry) Error(in interface{}) {
	s.GeneralLogger(in, logrus.ErrorLevel, false)
}
func (s *GoSentry) Warn(in interface{}) {
	s.GeneralLogger(in, logrus.WarnLevel, false)
}
func (s *GoSentry) Info(in interface{}) {
	s.GeneralLogger(in, logrus.InfoLevel, false)
}
func (s *GoSentry) Debug(in interface{}) {
	s.GeneralLogger(in, logrus.DebugLevel, false)
}
