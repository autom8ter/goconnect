package util

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/autom8ter/gosaas/sdk/go/proto/accounts"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	validator "gopkg.in/go-playground/validator.v9"
)

var validate = validator.New()

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password must not be empty")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return string(hash[:]), err
	}
	return string(hash[:]), nil
}

func ComparePasswordToHash(acc *accounts.Account, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(acc.HashedPassword), []byte(password))
}

func Validate(obj interface{}) error {
	return validate.Struct(obj)
}

func AsMap(obj interface{}) map[string]interface{} {
	struc := structs.New(obj)
	return struc.Map()
}
