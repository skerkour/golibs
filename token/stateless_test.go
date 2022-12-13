package token_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/skerkour/golibs/crypto"
	"github.com/skerkour/golibs/token"
	"github.com/skerkour/golibs/ulid"
)

func TestParseStateless(t *testing.T) {
	key, err := crypto.NewAEADKey()
	if err != nil {
		t.Errorf("Generating Key")
	}

	for i := 0; i < 2000; i += 1 {
		data := strconv.Itoa(i)
		newToken, err := token.NewStateless(key, ulid.New(), time.Now().Add(24*time.Hour), data)
		if err != nil {
			t.Errorf("Generating stateless token: %v", err)
		}
		parsedToken, err := token.ParseStateless(newToken.String())
		if err != nil {
			t.Errorf("parsing  statelesstoken: %v", err)
		}

		if parsedToken.Version() != newToken.Version() {
			t.Errorf("token.Version (%v) != parsedToken.Version (%v)", newToken.Version(), parsedToken.Version())
		}

		if !parsedToken.ID().Equal(newToken.ID()) {
			t.Errorf("token.ID (%v) != parsedToken.ID (%v)", newToken.ID().String(), parsedToken.ID().String())
		}

		if parsedToken.String() != newToken.String() {
			t.Errorf("token.String() (%v) != parsedToken.String() (%v)", newToken.String(), parsedToken.String())
		}

		if parsedToken.Data() != newToken.Data() {
			t.Errorf("token.Data() (%v) != parsedToken.Data() (%v)", newToken.Data(), parsedToken.Data())
		}

		if parsedToken.Data() != data {
			t.Errorf("token.Data() (%v) != data (%v)", newToken.Data(), data)
		}
	}
}

func TestVerifyStateless(t *testing.T) {
	wrongKey, err := crypto.NewAEADKey()
	if err != nil {
		t.Errorf("Generating wrong Key")
	}

	for i := 0; i < 2000; i += 1 {
		key, err := crypto.NewAEADKey()
		if err != nil {
			t.Errorf("Generating Key")
		}

		data := strconv.Itoa(i)
		newToken, err := token.NewStateless(key, ulid.New(), time.Now().Add(24*time.Hour), data)
		if err != nil {
			t.Errorf("Generating stateless token: %v", err)
		}

		err = newToken.Verify(key)
		if err != nil {
			t.Errorf("Verifying stateless token: %v", err)
		}

		err = newToken.Verify(wrongKey)
		if err == nil {
			t.Errorf("Accepting wrong key: %v", err)
		}
	}
}
