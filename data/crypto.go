package data

import (
	"crypto/rand"
	"encoding/base64"
	goerr "errors"

	"github.com/clockworksoul/cog2/errors"
	"golang.org/x/crypto/bcrypt"
)

var ErrCryptoHash = goerr.New("failed to generate password hash")

var ErrCryptoIO = goerr.New("failed to retieve randomness")

// CompareHashAndPassword receives a plaintext password and its hash, and
// returns true if they match.
func CompareHashAndPassword(hashedPassword string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

// GenerateRandomToken generates a random character token.
func GenerateRandomToken(length int) (string, error) {
	byteCount := (length * 3) / 4
	bytes := make([]byte, byteCount)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", errors.Wrap(ErrCryptoIO, err)
	}

	sEnc := base64.StdEncoding.EncodeToString(bytes)

	return sEnc, nil
}

// HashPassword receives a plaintext password and returns its hashed equivalent.
func HashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", errors.Wrap(ErrCryptoHash, err)
	}

	return string(hash), nil
}
