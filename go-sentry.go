package gosentry

import (
	"strings"
	"time"

	sentry "github.com/getsentry/sentry-go"
	logrus "github.com/sirupsen/logrus"
)

type GoSentry struct {
}

func NewSession(sentry_dsn string) (*GoSentry, error) {
	s := &GoSentry{}
	err := s.Init(sentry_dsn)
	return s, err
}

func SetLevel(level logrus.Level) {
	logrus.SetLevel(level)
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

func (s *GoSentry) GeneralLogger(in interface{}, level logrus.Level) {
	var e error
	var msg string

	// Extract the message from the interface
	switch v := in.(type) {
	case string:
		msg = v
	case error:
		msg = v.Error()
	default:
		msg = "SimpleSentry Panic: log only accepts 'string' or 'error' types"
		logrus.Fatal(msg)
		panic(msg)
	}

	// Switch on the log level
	switch level {
	case logrus.PanicLevel:
		sentry.CaptureException(e)
	case logrus.FatalLevel:
		sentry.CaptureException(e)
	case logrus.ErrorLevel:
		sentry.CaptureException(e)
	case logrus.WarnLevel:
		sentry.CaptureMessage(msg)
	case logrus.InfoLevel:
		sentry.CaptureMessage(msg)
	case logrus.DebugLevel:
		sentry.CaptureMessage(msg)
	}
}

func (s *GoSentry) Panic(in interface{}) {
	s.GeneralLogger(in, logrus.PanicLevel)
}
func (s *GoSentry) Fatal(in interface{}) {
	s.GeneralLogger(in, logrus.FatalLevel)

}
func (s *GoSentry) Error(in interface{}) {
	s.GeneralLogger(in, logrus.ErrorLevel)

}
func (s *GoSentry) Warn(in interface{}) {
	s.GeneralLogger(in, logrus.WarnLevel)

}
func (s *GoSentry) Info(in interface{}) {
	s.GeneralLogger(in, logrus.InfoLevel)

}
func (s *GoSentry) Debug(in interface{}) {
	s.GeneralLogger(in, logrus.DebugLevel)
}
