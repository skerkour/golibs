package rz

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(Writer(out), Fields(Timestamp(false)))
		log.Log("")
		if got, want := decodeIfBinaryToString(out.Bytes()), "{}\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("one-field", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(Writer(out), Fields(Timestamp(false)))
		log.Log("", String("foo", "bar"))
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"foo":"bar"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("two-field", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(Writer(out), Fields(Timestamp(false)))
		log.Log("", String("foo", "bar"), Int("n", 123))
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"foo":"bar","n":123}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
}

func TestInfo(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(Writer(out), Fields(Timestamp(false)))
		log.Info("")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("one-field", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(Writer(out), Fields(Timestamp(false)))
		log.Info("", String("foo", "bar"))
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","foo":"bar"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("two-field", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(Writer(out), Fields(Timestamp(false)))
		log.Info("", String("foo", "bar"), Int("n", 123))
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","foo":"bar","n":123}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
}

func TestFields(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(
		Writer(out),
		Fields(
			Timestamp(false),
			String("string", "foo"),
			Bytes("bytes", []byte("bar")),
			Hex("hex", []byte{0x12, 0xef}),
			RawJSON("json", []byte(`{"some":"json"}`)),
			Error("some_err", nil),
			Err(errors.New("some error")),
			Bool("bool", true),
			Int("int", 1),
			Int8("int8", 2),
			Int16("int16", 3),
			Int32("int32", 4),
			Int64("int64", 5),
			Uint("uint", 6),
			Uint8("uint8", 7),
			Uint16("uint16", 8),
			Uint32("uint32", 9),
			Uint64("uint64", 10),
			Float32("float32", 11.101),
			Float64("float64", 12.30303),
			Time("time", time.Time{}),
			Caller(true),
		),
	)
	_, file, line, _ := runtime.Caller(0)
	caller := fmt.Sprintf("%s:%d", file, line+2)
	log.Log("")
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"string":"foo","bytes":"bar","hex":"12ef","json":{"some":"json"},"error":"some error","bool":true,"int":1,"int8":2,"int16":3,"int32":4,"int64":5,"uint":6,"uint8":7,"uint16":8,"uint32":9,"uint64":10,"float32":11.101,"float64":12.30303,"time":"0001-01-01T00:00:00Z","caller":"`+caller+`"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestFieldsMap(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(Writer(out), Fields(Timestamp(false)))
	log.Log("", Map(map[string]interface{}{
		"nil":     nil,
		"string":  "foo",
		"bytes":   []byte("bar"),
		"error":   errors.New("some error"),
		"errors":  []error{errors.New("some error"), errors.New("some other error")},
		"bool":    true,
		"int":     int(1),
		"int8":    int8(2),
		"int16":   int16(3),
		"int32":   int32(4),
		"int64":   int64(5),
		"uint":    uint(6),
		"uint8":   uint8(7),
		"uint16":  uint16(8),
		"uint32":  uint32(9),
		"uint64":  uint64(10),
		"float32": float32(11),
		"float64": float64(12),
		"ipv6":    net.IP{0x20, 0x01, 0x0d, 0xb8, 0x85, 0xa3, 0x00, 0x00, 0x00, 0x00, 0x8a, 0x2e, 0x03, 0x70, 0x73, 0x34},
		"dur":     1 * time.Second,
		"time":    time.Time{},
		"obj":     obj{"a", "b", 1},
	}),
	)
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"bool":true,"bytes":"bar","dur":1000,"error":"some error","errors":["some error","some other error"],"float32":11,"float64":12,"int":1,"int16":3,"int32":4,"int64":5,"int8":2,"ipv6":"2001:db8:85a3::8a2e:370:7334","nil":null,"obj":{"Pub":"a","Tag":"b","priv":1},"string":"foo","time":"0001-01-01T00:00:00Z","uint":6,"uint16":8,"uint32":9,"uint64":10,"uint8":7}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestFieldsMapPnt(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(Writer(out), Fields(Timestamp(false)))
	log.Log("", Map(map[string]interface{}{
		"string":  new(string),
		"bool":    new(bool),
		"int":     new(int),
		"int8":    new(int8),
		"int16":   new(int16),
		"int32":   new(int32),
		"int64":   new(int64),
		"uint":    new(uint),
		"uint8":   new(uint8),
		"uint16":  new(uint16),
		"uint32":  new(uint32),
		"uint64":  new(uint64),
		"float32": new(float32),
		"float64": new(float64),
		"dur":     new(time.Duration),
		"time":    new(time.Time),
	}),
	)
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"bool":false,"dur":0,"float32":0,"float64":0,"int":0,"int16":0,"int32":0,"int64":0,"int8":0,"string":"","time":"0001-01-01T00:00:00Z","uint":0,"uint16":0,"uint32":0,"uint64":0,"uint8":0}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestFieldsMapNilPnt(t *testing.T) {
	var (
		stringPnt  *string
		boolPnt    *bool
		intPnt     *int
		int8Pnt    *int8
		int16Pnt   *int16
		int32Pnt   *int32
		int64Pnt   *int64
		uintPnt    *uint
		uint8Pnt   *uint8
		uint16Pnt  *uint16
		uint32Pnt  *uint32
		uint64Pnt  *uint64
		float32Pnt *float32
		float64Pnt *float64
		durPnt     *time.Duration
		timePnt    *time.Time
	)
	out := &bytes.Buffer{}
	log := New(Writer(out), Fields(Timestamp(false)))
	fields := map[string]interface{}{
		"string":  stringPnt,
		"bool":    boolPnt,
		"int":     intPnt,
		"int8":    int8Pnt,
		"int16":   int16Pnt,
		"int32":   int32Pnt,
		"int64":   int64Pnt,
		"uint":    uintPnt,
		"uint8":   uint8Pnt,
		"uint16":  uint16Pnt,
		"uint32":  uint32Pnt,
		"uint64":  uint64Pnt,
		"float32": float32Pnt,
		"float64": float64Pnt,
		"dur":     durPnt,
		"time":    timePnt,
	}
	log.Log("", Map(fields))
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"bool":null,"dur":null,"float32":null,"float64":null,"int":null,"int16":null,"int32":null,"int64":null,"int8":null,"string":null,"time":null,"uint":null,"uint16":null,"uint32":null,"uint64":null,"uint8":null}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestFields2(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(Writer(out), Fields(Timestamp(false)))
	_, file, line, _ := runtime.Caller(0)
	caller := fmt.Sprintf("%s:%d", file, line+2)
	log.Log("", Caller(true),
		String("string", "foo"),
		Bytes("bytes", []byte("bar")),
		Hex("hex", []byte{0x12, 0xef}),
		RawJSON("json", []byte(`{"some":"json"}`)),
		Error("some_err", nil),
		Err(errors.New("some error")),
		Bool("bool", true),
		Int("int", 1),
		Int8("int8", 2),
		Int16("int16", 3),
		Int32("int32", 4),
		Int64("int64", 5),
		Uint("uint", 6),
		Uint8("uint8", 7),
		Uint16("uint16", 8),
		Uint32("uint32", 9),
		Uint64("uint64", 10),
		IP("IPv4", net.IP{192, 168, 0, 100}),
		IP("IPv6", net.IP{0x20, 0x01, 0x0d, 0xb8, 0x85, 0xa3, 0x00, 0x00, 0x00, 0x00, 0x8a, 0x2e, 0x03, 0x70, 0x73, 0x34}),
		HardwareAddr("Mac", net.HardwareAddr{0x00, 0x14, 0x22, 0x01, 0x23, 0x45}),
		IPNet("IPNet", net.IPNet{IP: net.IP{192, 168, 0, 100}, Mask: net.CIDRMask(24, 32)}),
		Float32("float32", 11.1234),
		Float64("float64", 12.321321321),
		Duration("dur", 1*time.Second),
		Time("time", time.Time{}),
	)
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"string":"foo","bytes":"bar","hex":"12ef","json":{"some":"json"},"error":"some error","bool":true,"int":1,"int8":2,"int16":3,"int32":4,"int64":5,"uint":6,"uint8":7,"uint16":8,"uint32":9,"uint64":10,"IPv4":"192.168.0.100","IPv6":"2001:db8:85a3::8a2e:370:7334","Mac":"00:14:22:01:23:45","IPNet":"192.168.0.100/24","float32":11.1234,"float64":12.321321321,"dur":1000,"time":"0001-01-01T00:00:00Z","caller":"`+caller+`"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestFieldsArrayEmpty(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(Writer(out), Fields(Timestamp(false)))
	log.Log("", Strings("string", []string{}),
		Errors("err", []error{}),
		Bools("bool", []bool{}),
		Ints("int", []int{}),
		Ints8("int8", []int8{}),
		Ints16("int16", []int16{}),
		Ints32("int32", []int32{}),
		Ints64("int64", []int64{}),
		Uints("uint", []uint{}),
		Uints8("uint8", []uint8{}),
		Uints16("uint16", []uint16{}),
		Uints32("uint32", []uint32{}),
		Uints64("uint64", []uint64{}),
		Floats32("float32", []float32{}),
		Floats64("float64", []float64{}),
		Durations("dur", []time.Duration{}),
		Times("time", []time.Time{}),
	)
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"string":[],"err":[],"bool":[],"int":[],"int8":[],"int16":[],"int32":[],"int64":[],"uint":[],"uint8":[],"uint16":[],"uint32":[],"uint64":[],"float32":[],"float64":[],"dur":[],"time":[]}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestFieldsArraySingleElement(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(Writer(out), Fields(Timestamp(false)))
	log.Log("", Strings("string", []string{"foo"}),
		Errors("err", []error{errors.New("some error")}),
		Bools("bool", []bool{true}),
		Ints("int", []int{1}),
		Ints8("int8", []int8{2}),
		Ints16("int16", []int16{3}),
		Ints32("int32", []int32{4}),
		Ints64("int64", []int64{5}),
		Uints("uint", []uint{6}),
		Uints8("uint8", []uint8{7}),
		Uints16("uint16", []uint16{8}),
		Uints32("uint32", []uint32{9}),
		Uints64("uint64", []uint64{10}),
		Floats32("float32", []float32{11}),
		Floats64("float64", []float64{12}),
		Durations("dur", []time.Duration{1 * time.Second}),
		Times("time", []time.Time{{}}),
	)
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"string":["foo"],"err":["some error"],"bool":[true],"int":[1],"int8":[2],"int16":[3],"int32":[4],"int64":[5],"uint":[6],"uint8":[7],"uint16":[8],"uint32":[9],"uint64":[10],"float32":[11],"float64":[12],"dur":[1000],"time":["0001-01-01T00:00:00Z"]}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestFieldsArrayMultipleElement(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(Writer(out), Fields(Timestamp(false)))
	log.Log("", Strings("string", []string{"foo", "bar"}),
		Errors("err", []error{errors.New("some error"), nil}),
		Bools("bool", []bool{true, false}),
		Ints("int", []int{1, 0}),
		Ints8("int8", []int8{2, 0}),
		Ints16("int16", []int16{3, 0}),
		Ints32("int32", []int32{4, 0}),
		Ints64("int64", []int64{5, 0}),
		Uints("uint", []uint{6, 0}),
		Uints8("uint8", []uint8{7, 0}),
		Uints16("uint16", []uint16{8, 0}),
		Uints32("uint32", []uint32{9, 0}),
		Uints64("uint64", []uint64{10, 0}),
		Floats32("float32", []float32{11, 0}),
		Floats64("float64", []float64{12, 0}),
		Durations("dur", []time.Duration{1 * time.Second, 0}),
		Times("time", []time.Time{{}, {}}),
	)
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"string":["foo","bar"],"err":["some error",null],"bool":[true,false],"int":[1,0],"int8":[2,0],"int16":[3,0],"int32":[4,0],"int64":[5,0],"uint":[6,0],"uint8":[7,0],"uint16":[8,0],"uint32":[9,0],"uint64":[10,0],"float32":[11,0],"float64":[12,0],"dur":[1000,0],"time":["0001-01-01T00:00:00Z","0001-01-01T00:00:00Z"]}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestFieldsDisabled(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(Writer(out), Fields(Timestamp(false)), Level(InfoLevel))
	log.Debug("", String("string", "foo"),
		Bytes("bytes", []byte("bar")),
		Hex("hex", []byte{0x12, 0xef}),
		Error("some_err", nil),
		Err(errors.New("some error")),
		Bool("bool", true),
		Int("int", 1),
		Int8("int8", 2),
		Int16("int16", 3),
		Int32("int32", 4),
		Int64("int64", 5),
		Uint("uint", 6),
		Uint8("uint8", 7),
		Uint16("uint16", 8),
		Uint32("uint32", 9),
		Uint64("uint64", 10),
		Float32("float32", 11),
		Float64("float64", 12),
		Duration("dur", 1*time.Second),
		Time("time", time.Time{}),
	)
	if got, want := decodeIfBinaryToString(out.Bytes()), ""; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestFieldsAndFieldsCombined(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(Writer(out), Fields(Timestamp(false)), Fields(String("f1", "val"), String("f2", "val")))
	log.Log("", String("f3", "val"))
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"f1":"val","f2":"val","f3":"val"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestLevel(t *testing.T) {
	t.Run("Disabled", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(Writer(out), Fields(Timestamp(false)), Level(Disabled))
		log.Info("test")
		if got, want := decodeIfBinaryToString(out.Bytes()), ""; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("NoLevel/Disabled", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(Writer(out), Fields(Timestamp(false)), Level(Disabled))
		log.Log("test")
		if got, want := decodeIfBinaryToString(out.Bytes()), ""; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("NoLevel/Info", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(Writer(out), Fields(Timestamp(false)), Level(InfoLevel))
		log.Log("test")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"message":"test"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("NoLevel/Panic", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(Writer(out), Fields(Timestamp(false)), Level(PanicLevel))
		log.Log("test")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"message":"test"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("NoLevel/Log", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(Writer(out), Fields(Timestamp(false)), Level(InfoLevel))
		log.Log("test")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"message":"test"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})

	t.Run("Info", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(Writer(out), Fields(Timestamp(false)), Level(InfoLevel))
		log.Info("test")
		if got, want := decodeIfBinaryToString(out.Bytes()), `{"level":"info","message":"test"}`+"\n"; got != want {
			t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
		}
	})
}

func TestSampling(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(Writer(out), Fields(Timestamp(false)), Sampler(&SamplerBasic{N: 2}))
	log.Log("", Int("i", 1))
	log.Log("", Int("i", 2))
	log.Log("", Int("i", 3))
	log.Log("", Int("i", 4))
	if got, want := decodeIfBinaryToString(out.Bytes()), "{\"i\":1}\n{\"i\":3}\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestDiscard(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(Writer(out), Fields(Timestamp(false)))
	log.Log("test123", String("a", "b"), Discard())
	if got, want := decodeIfBinaryToString(out.Bytes()), ""; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}

	// Double call
	log.Log("test123", Discard(), String("a", "b"), Discard())
	if got, want := decodeIfBinaryToString(out.Bytes()), ""; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

type levelWriter struct {
	ops []struct {
		l LogLevel
		p string
	}
}

func (lw *levelWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (lw *levelWriter) WriteLevel(lvl LogLevel, p []byte) (int, error) {
	p = decodeIfBinaryToBytes(p)
	lw.ops = append(lw.ops, struct {
		l LogLevel
		p string
	}{lvl, string(p)})
	return len(p), nil
}

func TestLevelWriter(t *testing.T) {
	lw := &levelWriter{
		ops: []struct {
			l LogLevel
			p string
		}{},
	}
	log := New(Writer(lw), Fields(Timestamp(false)))
	log.Debug("1")
	log.Info("2")
	log.Warn("3")
	log.Error("4")
	log.Log("nolevel-1")
	log.LogWithLevel(DebugLevel, "5")

	want := []struct {
		l LogLevel
		p string
	}{
		{DebugLevel, `{"level":"debug","message":"1"}` + "\n"},
		{InfoLevel, `{"level":"info","message":"2"}` + "\n"},
		{WarnLevel, `{"level":"warning","message":"3"}` + "\n"},
		{ErrorLevel, `{"level":"error","message":"4"}` + "\n"},
		{NoLevel, `{"message":"nolevel-1"}` + "\n"},
		{DebugLevel, `{"level":"debug","message":"5"}` + "\n"},
	}
	if got := lw.ops; !reflect.DeepEqual(got, want) {
		t.Errorf("invalid ops:\ngot:\n%v\nwant:\n%v", got, want)
	}
}

func TestContextTimestamp(t *testing.T) {
	tfn := func() time.Time {
		return time.Date(2001, time.February, 3, 4, 5, 6, 7, time.UTC)
	}
	out := &bytes.Buffer{}
	log := New(Writer(out), Fields(String("foo", "bar")), TimestampFunc(tfn))
	log.Log("hello world")

	if got, want := decodeIfBinaryToString(out.Bytes()), `{"foo":"bar","timestamp":"2001-02-03T04:05:06Z","message":"hello world"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

func TestEventTimestamp(t *testing.T) {
	tfn := func() time.Time {
		return time.Date(2001, time.February, 3, 4, 5, 6, 7, time.UTC)
	}
	out := &bytes.Buffer{}
	log := New(Writer(out), Fields(Timestamp(false), String("foo", "bar")), TimestampFunc(tfn))
	log.Log("hello world", Timestamp(true))

	if got, want := decodeIfBinaryToString(out.Bytes()), `{"foo":"bar","timestamp":"2001-02-03T04:05:06Z","message":"hello world"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

type loggableError struct {
	error
}

func (l loggableError) MarshalRzObject(e *Event) {
	e.Append(String("message", l.error.Error()+": loggableError"))
}

func TestErrorMarshalFunc(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(Writer(out), Fields(Timestamp(false)))

	// test default behavior
	log.Log("msg", Err(errors.New("err")))
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"error":"err","message":"msg"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
	out.Reset()

	log.Log("msg", Err(loggableError{errors.New("err")}))
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"error":{"message":"err: loggableError"},"message":"msg"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
	out.Reset()

	// test overriding the ErrorMarshalFunc
	originalErrorMarshalFunc := ErrorMarshalFunc
	defer func() {
		ErrorMarshalFunc = originalErrorMarshalFunc
	}()

	ErrorMarshalFunc = func(err error) interface{} {
		return err.Error() + ": marshaled string"
	}
	log.Log("msg", Err(errors.New("err")))
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"error":"err: marshaled string","message":"msg"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}

	out.Reset()
	ErrorMarshalFunc = func(err error) interface{} {
		return errors.New(err.Error() + ": new error")
	}
	log.Log("msg", Err(errors.New("err")))
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"error":"err: new error","message":"msg"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}

	out.Reset()
	ErrorMarshalFunc = func(err error) interface{} {
		return loggableError{err}
	}
	log.Log("msg", Err(errors.New("err")))
	if got, want := decodeIfBinaryToString(out.Bytes()), `{"error":{"message":"err: loggableError"},"message":"msg"}`+"\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}

type errWriter struct {
	error
}

func (w errWriter) Write(p []byte) (n int, err error) {
	return 0, w.error
}

func TestErrorHandler(t *testing.T) {
	var got error
	want := errors.New("write error")
	ErrorHandler = func(err error) {
		got = err
	}
	log := New(Writer(errWriter{want}))
	log.Log("test")
	if got != want {
		t.Errorf("ErrorHandler err = %#v, want %#v", got, want)
	}
}

func TestWrite(t *testing.T) {
	out := &bytes.Buffer{}
	log := New(Writer(out), Fields(Timestamp(false)))
	log.Write([]byte("test"))
	if got, want := decodeIfBinaryToString(out.Bytes()), "{\"message\":\"test\"}\n"; got != want {
		t.Errorf("invalid log output:\ngot:  %v\nwant: %v", got, want)
	}
}
