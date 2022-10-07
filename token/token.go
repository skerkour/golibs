package token

import (
	"strings"

	"github.com/skerkour/bloom2/domain/kernel"
	"github.com/skerkour/bloom2/errs"
	"github.com/skerkour/bloom2/libs/base32"
	"github.com/skerkour/bloom2/libs/crypto"
	"github.com/skerkour/bloom2/libs/ulid"
)

const (
	SecretSize = crypto.KeySize128
	HashSize   = crypto.KeySize512
)

var (
	ErrTokenIsNotValid = errs.InvalidArgument("Token is not valid.")
	ErrDataIsTooLong   = errs.InvalidArgument("data is too long.")
)

//TODO improve performance using arrays
type Token struct {
	id     ulid.ULID
	secret []byte
	hash   []byte
	str    string
}

func New() (token Token, err error) {
	secret, err := newSecret()
	if err != nil {
		return
	}

	return new("", secret)
}

func NewWithSecret(secret []byte) (token Token, err error) {
	return new("", secret)
}

func NewWithPrefix(prefix string) (token Token, err error) {
	secret, err := newSecret()
	if err != nil {
		return
	}
	return new(prefix, secret)
}

func newSecret() (secret []byte, err error) {
	secret, err = crypto.RandBytes(SecretSize)
	if err != nil {
		// TODO: log
		err = errs.Internal("token: Generating secret", err)
		return
	}
	return
}

func new(prefix string, secret []byte) (token Token, err error) {
	id := ulid.New()

	idBytes, _ := id.MarshalBinary()

	hash, err := crypto.DeriveKeyFromKey(secret, idBytes, HashSize)
	if err != nil {
		// TODO: log
		err = errs.Internal("token: Hashing secret", err)
		return
	}

	data := append(idBytes, secret...)
	str := base32.EncodeToString(data)
	str = prefix + str

	token = Token{
		id,
		secret,
		hash,
		str,
	}
	return
}

func (token *Token) String() string {
	return token.str
}

func (token *Token) ID() ulid.ULID {
	return token.id
}

func (token *Token) Secret() []byte {
	return token.secret
}

func (token *Token) Hash() []byte {
	return token.hash
}

func Parse(input string) (token Token, err error) {
	return ParseWithPrefix(input, "")
}

func ParseWithPrefix(input, prefix string) (token Token, err error) {
	var tokenBytes []byte

	token.str = input

	if prefix != "" {
		if !strings.HasPrefix(input, prefix) {
			err = kernel.ErrTokenIsNotValid
			return
		}
		input = strings.TrimPrefix(input, prefix)
	}

	tokenBytes, err = base32.DecodeString(input)
	if err != nil {
		err = kernel.ErrTokenIsNotValid
		return
	}

	if len(tokenBytes) != ulid.Size+SecretSize {
		err = kernel.ErrTokenIsNotValid
		return
	}

	tokenIDBytes := tokenBytes[:ulid.Size]
	token.secret = tokenBytes[ulid.Size:]

	token.id, err = ulid.ParseBytes(tokenIDBytes)
	if err != nil {
		err = kernel.ErrTokenIsNotValid
		return
	}

	token.hash, err = crypto.DeriveKeyFromKey(token.secret, tokenIDBytes, HashSize)
	if err != nil {
		// TODO: log
		err = errs.Internal("token: Hashing secret", err)
		return
	}

	return
}

func (token *Token) Verify(hash []byte) (err error) {
	// in case we need to update hash size later
	// if len(hash) == OldHashSize {
	// token.hash = crypto.DeriveKeyFromKey(secret, idBytes, OldHashSize)
	// ..
	// }

	if !crypto.ConstantTimeCompare(hash, token.hash) {
		err = kernel.ErrTokenIsNotValid
	}
	return
}
