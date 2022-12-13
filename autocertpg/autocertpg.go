package autocertpg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/skerkour/golibs/crypto"
	"github.com/skerkour/golibs/db"
	"golang.org/x/crypto/acme/autocert"
)

type Cache struct {
	key []byte
	db  db.DB
}

type cert struct {
	Key           string `db:"key"`
	EncryptedData []byte `db:"encrypted_data"`
}

func NewCache(db db.DB, key []byte) *Cache {
	return &Cache{
		db:  db,
		key: key,
	}
}

func (cache *Cache) Get(ctx context.Context, key string) (data []byte, err error) {
	var cert cert
	query := "SELECT * FROM certs WHERE key = $1"

	err = cache.db.Get(ctx, &cert, query, key)
	if err != nil {
		if err == sql.ErrNoRows {
			err = autocert.ErrCacheMiss
		}
		return
	}

	data, err = crypto.Decrypt(cache.key, cert.EncryptedData, []byte(cert.Key))
	if err != nil {
		err = fmt.Errorf("autocertpg: decrypting data: %w", err)
		return
	}

	return
}

func (cache *Cache) Put(ctx context.Context, key string, data []byte) (err error) {
	query := `
	INSERT INTO certs (key, encrypted_data)
		VALUES ($1, $2)
		ON CONFLICT (key)
		DO UPDATE SET encrypted_data = $2
	`

	encryptedData, err := crypto.Encrypt(cache.key, data, []byte(key))
	if err != nil {
		err = fmt.Errorf("autocertpg: encrypting data: %w", err)
		return
	}

	_, err = cache.db.Exec(ctx, query, key, encryptedData)
	if err != nil {
		return
	}

	return
}

func (cache *Cache) Delete(ctx context.Context, key string) (err error) {
	query := "DELETE FROM certs WHERE key = $1"

	_, err = cache.db.Exec(ctx, query, key)
	if err != nil {
		return
	}

	return
}
