package token

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/skerkour/libs/crypto"
	"github.com/skerkour/libs/ulid"
)

const maxStatelessData = 128

type Stateless struct {
	version   uint8
	payload   statelessPayload
	signature []byte
	str       string
}

type statelessPayload struct {
	ID   ulid.ULID `json:"id"`
	Exp  time.Time `json:"exp"`
	Data string    `json:"data"`
}

func (token *Stateless) String() string {
	return token.str
}

func (token *Stateless) Version() uint8 {
	return token.version
}

func (token *Stateless) ID() ulid.ULID {
	return token.payload.ID
}

func (token *Stateless) Data() string {
	return token.payload.Data
}

func (token *Stateless) Verify(key []byte) (err error) {
	parts := strings.Split(token.str, ".")
	if len(parts) != 3 {
		err = ErrTokenIsNotValid
		return
	}

	versionAndPayload := parts[0] + "." + parts[1]

	signature, err := crypto.Mac(key, []byte(versionAndPayload), crypto.KeySize256)
	if err != nil {
		return
	}

	if !crypto.ConstantTimeCompare(signature, token.signature) {
		err = ErrTokenIsNotValid
		return
	}

	return
}

func NewStateless(key []byte, id ulid.ULID, expire time.Time, data string) (token Stateless, err error) {
	if len(data) > maxStatelessData {
		err = ErrDataIsTooLong
		return
	}

	token.version = 1
	token.payload = statelessPayload{
		ID:   id,
		Exp:  expire.UTC(),
		Data: data,
	}

	payloadJson, err := json.Marshal(token.payload)
	if err != nil {
		return
	}

	payloadBase64 := base64.RawURLEncoding.EncodeToString(payloadJson)

	token.str = "v1" + "." + payloadBase64

	token.signature, err = crypto.Mac(key, []byte(token.str), crypto.KeySize256)
	if err != nil {
		return
	}

	token.str += "." + base64.RawURLEncoding.EncodeToString(token.signature)

	return
}

func ParseStateless(tokenStr string) (token Stateless, err error) {
	if len(tokenStr) > 170+maxStatelessData {
		err = ErrDataIsTooLong
		return
	}

	token.str = tokenStr
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		err = ErrTokenIsNotValid
		return
	}

	// Version
	switch parts[0] {
	case "v1":
		token.version = 1
	default:
		err = ErrTokenIsNotValid
		return
	}

	// Payload
	payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		err = ErrTokenIsNotValid
		return
	}

	err = json.Unmarshal(payloadJSON, &token.payload)
	if err != nil {
		err = ErrTokenIsNotValid
		return
	}

	// Signature
	token.signature, err = base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		err = ErrTokenIsNotValid
		return
	}

	return
}
