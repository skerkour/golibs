package rz

import (
	"reflect"
	"testing"
)

func TestEvent_Fields(t *testing.T) {
	fields := map[string]interface{}{
		"hostname": "localhost",
		"latency":  3000.0,
	}
	event := newEvent(nil, DebugLevel)

	for key, value := range fields {
		event.Append(Any(key, value))
	}

	got, err := event.Fields()
	if err != nil {
		t.Fatalf("Event.Fields() returned error: %s", err)
	}
	if !reflect.DeepEqual(fields, got) {
		t.Errorf("Event.Fields() returned %+v, want %+v", got, fields)
	}
}
