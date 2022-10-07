package rz

import (
	"context"
	"io/ioutil"
	"testing"
)

func TestFromCtx(t *testing.T) {
	log := New(Writer(ioutil.Discard))
	ctx := log.ToCtx(context.Background())
	log2 := FromCtx(ctx)
	if &log != log2 {
		t.Error("FromCtx did not return the expected logger")
	}

	// update
	log = log.With(Level(InfoLevel))
	ctx = log.ToCtx(ctx)
	log2 = FromCtx(ctx)
	if &log != log2 {
		t.Error("FromCtx did not return the expected logger")
	}

	log2 = FromCtx(context.Background())
	if log2 == nil || log2 == &log {
		t.Error("FromCtx did not return the expected logger")
	}
}

func TestFromCtxDisabled(t *testing.T) {
	dl := New(Writer(ioutil.Discard), Level(Disabled))
	ctx := dl.ToCtx(context.Background())
	if ctx != context.Background() {
		t.Error("ToCtx stored a disabled logger")
	}

	l := New(
		Writer(ioutil.Discard),
		Fields(String("foo", "bar")),
	)
	ctx = l.ToCtx(ctx)
	if FromCtx(ctx) != &l {
		t.Error("WithContext did not store logger")
	}

	l = l.With(Level(DebugLevel))
	ctx = l.ToCtx(ctx)
	if FromCtx(ctx) != &l {
		t.Error("ToCtx did not store copied logger")
	}

	ctx = dl.ToCtx(ctx)
	if FromCtx(ctx) != &dl {
		t.Error("ToCtx did not overide logger with a disabled logger")
	}
}
