// Package log provides a global logger for  rz.
package log

import (
	"github.com/skerkour/golibs/rz"
)

// logger is the global logger.
var logger = rz.New(rz.CallerSkipFrameCount(4))

// SetLogger update log's logger
func SetLogger(log rz.Logger) {
	logger = log.With(rz.CallerSkipFrameCount(4))
}

// Logger returns log's logger
func Logger() rz.Logger {
	return logger.With(rz.CallerSkipFrameCount(3))
}

// With duplicates the global logger and update it's configuration.
func With(options ...rz.LoggerOption) rz.Logger {
	options = append([]rz.LoggerOption{rz.CallerSkipFrameCount(3)}, options...)
	return logger.With(options...)
}

// LogWithLevel logs a new message with the given level.
func LogWithLevel(level rz.LogLevel, message string, fields ...rz.Field) {
	logger.LogWithLevel(level, message, fields...)
}

// Debug starts a new message with debug level.
func Debug(message string, fields ...rz.Field) {
	logger.Debug(message, fields...)
}

// Info logs a new message with info level.
func Info(message string, fields ...rz.Field) {
	logger.Info(message, fields...)
}

// Warn logs a new message with warn level.
func Warn(message string, fields ...rz.Field) {
	logger.Warn(message, fields...)
}

// Error logs a message with error level.
func Error(message string, fields ...rz.Field) {
	logger.Error(message, fields...)
}

// Fatal logs a new message with fatal level. The os.Exit(1) function
// is then called, which terminates the program immediately.
func Fatal(message string, fields ...rz.Field) {
	logger.Fatal(message, fields...)
}

// Panic logs a new message with panic level. The panic() function
// is then called, which stops the ordinary flow of a goroutine.
func Panic(message string, fields ...rz.Field) {
	logger.Panic(message, fields...)
}

// Log logs a new message with no level. Setting GlobalLevel to Disabled
// will still disable events produced by this method.
func Log(message string, fields ...rz.Field) {
	logger.Log(message, fields...)
}

// Append the fields to the internal logger's context.
// It does not create a new copy of the logger and rely on a mutex to enable thread safety,
// so `Config(With(fields...))` often is preferable.
func Append(fields ...rz.Field) {
	logger.Append(fields...)
}

// NewDict create a new Dict with the logger's configuration
func NewDict(fields ...rz.Field) *rz.Event {
	return logger.NewDict(fields...)
}
