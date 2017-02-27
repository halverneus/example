package encrypt

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/scrypt"
)

// NewBase64Salt creates a random salt encoded as Base64.
func NewBase64Salt() (salt string, err error) {
	// Read 24 random bytes.
	out := make([]byte, 24)
	var n int
	if n, err = io.ReadFull(rand.Reader, out); nil != err {
		return
	}

	// Verify successful reading of random data.
	if len(out) != n {
		err = errors.New("Count mismatch while creating salt")
		return
	}

	// Return encoded into base64.
	salt = base64.StdEncoding.EncodeToString(out)
	return
}

// SaltFromBase64 converts a Base64 encoded salt to bytes.
func SaltFromBase64(enc string) (salt []byte, err error) {
	return base64.StdEncoding.DecodeString(enc)
}

// Password gets encrypted here.
func Password(password, salt string) (enc string, err error) {
	// Convert password directly to byte slice and decode salt.
	passwordRaw := []byte(password)
	var saltRaw []byte
	if saltRaw, err = SaltFromBase64(salt); nil != err {
		return
	}

	// Encrypt password to byte slice. Encryption values determined by current
	// cryptography suggestion.
	var raw []byte
	if raw, err = scrypt.Key(passwordRaw, saltRaw, 16384, 8, 1, 32); nil != err {
		return
	}

	// Encode as Base64.
	enc = base64.StdEncoding.EncodeToString(raw)
	return
}
