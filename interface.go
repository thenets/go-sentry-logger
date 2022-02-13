package gosentry

import logrus "github.com/sirupsen/logrus"

var global_go_sentry *GoSentry

func NewSession(sentry_dsn string) (*GoSentry, error) {
	s := &GoSentry{}
	err := s.Init(sentry_dsn)
	global_go_sentry = s
	return s, err
}

func SetLevel(level logrus.Level) {
	logrus.SetLevel(level)
}

func Panic(in interface{}) {
	global_go_sentry.Panic(in)
}
func Fatal(in interface{}) {
	global_go_sentry.Fatal(in)
}
func Error(in interface{}) {
	global_go_sentry.Error(in)
}
func Warn(in interface{}) {
	global_go_sentry.Warn(in)
}
func Info(in interface{}) {
	global_go_sentry.Info(in)
}
func Debug(in interface{}) {
	global_go_sentry.Debug(in)
}
