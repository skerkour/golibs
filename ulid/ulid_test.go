package ulid_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/skerkour/golibs/ulid"
)

func TestNew(t *testing.T) {
	for i := 0; i < 1000; i += 1 {
		id := ulid.New()
		time.Sleep(time.Millisecond)
		fmt.Println(id.String(), "->", hex.EncodeToString(id[:]))
	}
}

func TestParse(t *testing.T) {
	for i := 0; i < 1000; i += 1 {
		id := ulid.New()
		parsed, err := ulid.Parse(id.String())
		if err != nil {
			t.Errorf("parsing ulid: %s", err)
		}
		if !bytes.Equal(id[:], parsed[:]) {
			t.Errorf("parsed (%s) != original ULID (%s)", parsed.String(), id.String())
		}
	}
}
