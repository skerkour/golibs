<p align="center">
  <h3 align="center">rz</h3>
  <p align="center">RipZap - Fast and 0 allocs leveled JSON logger for Go ⚡️. Dependency free. </p>
</p>

--------

[Make logging great again](https://kerkour.com/blog/logging/)

[![Go Reference](https://pkg.go.dev/badge/github.com/skerkour/libs/rz.svg)](https://pkg.go.dev/github.com/skerkour/libs/rz)
<!-- [![GitHub release](https://img.shields.io/github/release/skerkour/rz.svg)](https://github.com/skerkour/libs/rz/releases) -->
<!-- [![Coverage](http://gocover.io/_badge/github.com/skerkour/libs/rz)](http://gocover.io/github.com/skerkour/libs/rz) -->

![Console logging](docs/example_screenshot.png)

The rz package provides a fast and simple logger dedicated to JSON output avoiding allocations and reflection..

Uber's [zap](https://godoc.org/go.uber.org/zap) and rs's [zerolog](https://godoc.org/github.com/rs/zerolog)
libraries pioneered this approach.

ripzap is an update of zerolog taking this concept to the next level with a **simpler** to use and **safer**
API and even better [performance](#benchmarks).

To keep the code base and the API simple, ripzap focuses on efficient structured logging only.
Pretty logging on the console is made possible using the provided (but inefficient)
[`Formatter`s](https://godoc.org/github.com/skerkour/libs/rz#LogFormatter).


# Project status

rz is stable and the API is fixed. The project is considered as **done** and will not accept contributions other than for bug fixes.

Forks are welcome if you want to more features.


1. [Quickstart](#quickstart)
2. [Configuration](#configuration)
3. [Field types](#field-types)
4. [HTTP Handler](#http-handler)
5. [Examples](#examples)
6. [Benchmarks](#benchmarks)
7. [Contributing](#contributing)
8. [License](#license)

-------------------

## Quickstart

`rz` requires [Go modules](https://blog.golang.org/using-go-modules) so you need a `go.mod` file at the
root of your project.

```go
package main

import (
	"os"

	"github.com/skerkour/libs/rz"
	"github.com/skerkour/libs/rz/log"
)

func main() {

	env := os.Getenv("GO_ENV")
	hostname, _ := os.Hostname()

	// update global logger's context fields
	log.SetLogger(log.With(rz.Fields(
		rz.String("hostname", hostname), rz.String("environment", env),
	)))

	if env == "production" {
		log.SetLogger(log.With(rz.Level(rz.InfoLevel)))
	}

	log.Info("info from logger", rz.String("hello", "world"))
	// {"level":"info","hostname":"","environment":"","hello":"world","timestamp":"2019-02-07T09:30:07Z","message":"info from logger"}
}
```


## Configuration

### Logger
```go
// Writer update logger's writer.
func Writer(writer io.Writer) LoggerOption {}
// Level update logger's level.
func Level(lvl LogLevel) LoggerOption {}
// Sampler update logger's sampler.
func Sampler(sampler LogSampler) LoggerOption {}
// AddHook appends hook to logger's hook
func AddHook(hook LogHook) LoggerOption {}
// Hooks replaces logger's hooks
func Hooks(hooks ...LogHook) LoggerOption {}
// With replaces logger's context fields
func With(fields func(*Event)) LoggerOption {}
// Stack enable/disable stack in error messages.
func Stack(enableStack bool) LoggerOption {}
// Timestamp enable/disable timestamp logging in error messages.
func Timestamp(enableTimestamp bool) LoggerOption {}
// Caller enable/disable caller field in message messages.
func Caller(enableCaller bool) LoggerOption {}
// Formatter update logger's formatter.
func Formatter(formatter LogFormatter) LoggerOption {}
// TimestampFieldName update logger's timestampFieldName.
func TimestampFieldName(timestampFieldName string) LoggerOption {}
// LevelFieldName update logger's levelFieldName.
func LevelFieldName(levelFieldName string) LoggerOption {}
// MessageFieldName update logger's messageFieldName.
func MessageFieldName(messageFieldName string) LoggerOption {}
// ErrorFieldName update logger's errorFieldName.
func ErrorFieldName(errorFieldName string) LoggerOption {}
// CallerFieldName update logger's callerFieldName.
func CallerFieldName(callerFieldName string) LoggerOption {}
// CallerSkipFrameCount update logger's callerSkipFrameCount.
func CallerSkipFrameCount(callerSkipFrameCount int) LoggerOption {}
// ErrorStackFieldName update logger's errorStackFieldName.
func ErrorStackFieldName(errorStackFieldName string) LoggerOption {}
// TimeFieldFormat update logger's timeFieldFormat.
func TimeFieldFormat(timeFieldFormat string) LoggerOption {}
// TimestampFunc update logger's timestampFunc.
func TimestampFunc(timestampFunc func() time.Time) LoggerOption {}
```

### Global
```go
var (
	// DurationFieldUnit defines the unit for time.Duration type fields added
	// using the Duration method.
	DurationFieldUnit = time.Millisecond

	// DurationFieldInteger renders Duration fields as integer instead of float if
	// set to true.
	DurationFieldInteger = false

	// ErrorHandler is called whenever rz fails to write an event on its
	// output. If not set, an error is printed on the stderr. This handler must
	// be thread safe and non-blocking.
	ErrorHandler func(err error)
)
```


## Field Types

### Standard Types

* `String`
* `Bool`
* `Int`, `Int8`, `Int16`, `Int32`, `Int64`
* `Uint`, `Uint8`, `Uint16`, `Uint32`, `Uint64`
* `Float32`, `Float64`

### Advanced Fields

* `Err`: Takes an `error` and render it as a string using the `logger.errorFieldName` field name.
* `Error`: Adds a field with a `error`.
* `Timestamp`: Insert a timestamp field with `logger.timestampFieldName` field name and formatted using `logger.timeFieldFormat`.
* `Time`: Adds a field with the time formated with the `logger.timeFieldFormat`.
* `Duration`: Adds a field with a `time.Duration`.
* `Dict`: Adds a sub-key/value as a field of the event.
* `Interface`: Uses reflection to marshal the type.


## HTTP Handler

See the [skerkour/rz/rzhttp](https://godoc.org/github.com/skerkour/libs/rz/rzhttp) package or the
[example here](https://github.com/skerkour/libs/rz/tree/master/examples/http).


## Examples

See the [examples](https://github.com/skerkour/libs/rz/tree/master/examples) folder.


## Benchmarks

See [Logbench](http://hackemist.com/logbench/)

or

```
$ make benchmarks
```

## Contributing

See [https://bloom.sh/contribute](https://bloom.sh/contribute)


## License

`Apache 2.0`. See `LICENSE.txt`

From an original work by [rs](https://github.com/rs): [zerolog](https://github.com/rs/zerolog) - commit aa55558e4cb2e8f05cd079342d430f77e946d00a
