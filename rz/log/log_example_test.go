package log_test

import (
	"errors"
	"flag"
	"fmt"
	"time"

	"github.com/skerkour/golibs/rz"
	"github.com/skerkour/golibs/rz/log"
)

// setup would normally be an init() function, however, there seems
// to be something awry with the testing framework when we set the
// global Logger from an init()
func setup() {
	// In order to always output a static time to stdout for these
	// examples to pass, we need to override rz.TimestampFunc
	// and log.Logger globals -- you would not normally need to do this
	log.SetLogger(rz.New(rz.TimeFieldFormat(""), rz.TimestampFunc(func() time.Time {
		return time.Date(2008, 1, 8, 17, 5, 05, 0, time.UTC)
	})))
}

// Example of a log with no particular "level"
func ExampleLog() {
	setup()
	log.Log("hello world")

	// Output: {"timestamp":1199811905,"message":"hello world"}
}

// Example of a log at a particular "level" (in this case, "debug")
func ExampleDebug() {
	setup()
	log.Debug("hello world")

	// Output: {"level":"debug","timestamp":1199811905,"message":"hello world"}
}

// Example of a log at a particular "level" (in this case, "info")
func ExampleInfo() {
	setup()
	log.Info("hello world")

	// Output: {"level":"info","timestamp":1199811905,"message":"hello world"}
}

// Example of a log at a particular "level" (in this case, "warn")
func ExampleWarn() {
	setup()
	log.Warn("hello world")

	// Output: {"level":"warning","timestamp":1199811905,"message":"hello world"}
}

// Example of a log at a particular "level" (in this case, "error")
func ExampleError() {
	setup()
	log.Error("hello world")

	// Output: {"level":"error","timestamp":1199811905,"message":"hello world"}
}

// Example of a log at a particular "level" (in this case, "fatal")
func ExampleFatal() {
	setup()
	err := errors.New("A repo man spends his life getting into tense situations")
	service := "myservice"

	log.Fatal(fmt.Sprintf("Cannot start %s", service), rz.Err(err), rz.String("service", service))

	// Outputs: {"level":"fatal","timestamp":1199811905,"error":"A repo man spends his life getting into tense situations","service":"myservice","message":"Cannot start myservice"}
}

// TODO: Panic

// This example uses command-line flags to demonstrate various outputs
// depending on the chosen log level.
func Example() {
	setup()
	debug := flag.Bool("debug", false, "sets log level to debug")

	flag.Parse()

	// Default level for this example is info, unless debug flag is present

	if *debug {
		logger := log.Logger()
		defer func() {
			log.SetLogger(logger)
		}()
		log.SetLogger(log.With(rz.Level(rz.DebugLevel)))
	}

	log.Info("This message appears when log level set to Debug or Info")

	// Output: {"level":"info","timestamp":1199811905,"message":"This message appears when log level set to Debug or Info"}
}

// TODO: Output

// TODO: With

// TODO: Level

// TODO: Sample

// TODO: Hook

// TODO: WithLevel

// TODO: Ctx
