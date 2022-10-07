package ulid

import (
	"bytes"
	"crypto/rand"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/skerkour/libs/base32"
)

// TODO
type ULID [Size]byte

const (
	Size       = 16
	StringSize = 26
)

var (
	Nil ULID
)

func New() ULID {
	ulid, err := NewErr()
	if err != nil {
		panic(err)
	}
	return ulid
}

func NewErr() (ULID, error) {
	var ulid ULID

	err := ulid.setRandom()
	if err != nil {
		return Nil, err
	}

	ulid.setTime(time.Now())
	return ulid, nil
}

func (ulid *ULID) setTime(t time.Time) {
	timestamp := uint64(t.UnixNano() / int64(time.Millisecond))
	(*ulid)[0] = byte(timestamp >> 40)
	(*ulid)[1] = byte(timestamp >> 32)
	(*ulid)[2] = byte(timestamp >> 24)
	(*ulid)[3] = byte(timestamp >> 16)
	(*ulid)[4] = byte(timestamp >> 8)
	(*ulid)[5] = byte(timestamp)
	// var x, y byte
	// timestamp := uint64(t.UnixNano() / int64(time.Millisecond))
	// // Backups [6] and [7] bytes to override them with their original values later.
	// x, y, ulid[6], ulid[7] = ulid[6], ulid[7], x, y
	// binary.LittleEndian.PutUint64(ulid[:], timestamp)
	// // Truncates at the 6th byte as designed in the original spec (48 bytes).
	// ulid[6], ulid[7] = x, y
}

func (ulid *ULID) setRandom() (err error) {
	_, err = rand.Read(ulid[6:])
	return
}

func Parse(input string) (ret ULID, err error) {
	var retBytes []byte

	retBytes, err = base32.DecodeString(input)
	if err != nil {
		return
	}
	if len(retBytes) != Size {
		err = fmt.Errorf("invalid ULID (got %d bytes)", len(input))
		return
	}

	copy(ret[:], retBytes)
	return
}

func ParseBytes(input []byte) (ret ULID, err error) {
	switch len(input) {
	case Size:
		copy(ret[:], input)
	case StringSize:
		ret, err = Parse(string(input))
	default:
		err = fmt.Errorf("invalid ULID (got %d bytes)", len(input))
	}

	return
}

// MarshalText implements encoding.TextMarshaler.
func (ulid ULID) MarshalText() ([]byte, error) {
	str := ulid.String()
	return []byte(str), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (ulid *ULID) UnmarshalText(data []byte) error {
	id, err := ParseBytes(data)
	if err != nil {
		return err
	}
	*ulid = id
	return nil
}

// String returns the string form of ulid
func (ulid ULID) String() string {
	return base32.EncodeToString(ulid[:])
}

// UUIDString returns the UUID string form of ulid, xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
func (ulid ULID) UUIDString() string {
	uuid := uuid.UUID(ulid)
	return uuid.String()
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (ulid ULID) MarshalBinary() ([]byte, error) {
	return ulid[:], nil
}

func (ulid ULID) Bytes() []byte {
	return ulid[:]
}

// TODO: improve
// Scan implements sql.Scanner so ULIDs can be read from databases transparently.
// Currently, database types that map to string and []byte are supported. Please
// consult database-specific driver documentation for matching types.
func (ulid *ULID) Scan(src interface{}) (err error) {
	var uuid uuid.UUID
	err = uuid.Scan(src)
	*ulid = ULID(uuid)
	return
}

// Value implements sql.Valuer so that ULIDs can be written to databases
// transparently. Currently, ULIDs map to strings. Please consult
// database-specific driver documentation for matching types.
func (ulid ULID) Value() (driver.Value, error) {
	return ulid.UUIDString(), nil
}

func (ulid ULID) Equal(other ULID) bool {
	return bytes.Equal(ulid[:], ulid[:])
}
