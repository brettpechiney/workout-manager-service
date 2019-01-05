package logging

// IshiLogger represents a logger used in all Ishi golang projects. It is
// backed by the zap sugared logger.
type IshiLogger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
	Panic(msg string, args ...interface{})
	Sync() error
	WithFields(args ...interface{}) IshiLogger
}
