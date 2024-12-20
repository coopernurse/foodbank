package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"strconv"
	"strings"
	"time"
)

// EncryptSessionToken encrypts the person ID and expiry time into a session token.
func EncryptSessionToken(personID string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	expiry := time.Now().Add(24 * time.Hour).Unix()
	plaintext := []byte(personID + "|" + strconv.FormatInt(expiry, 10))

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptSessionToken decrypts the session token back to the person ID and expiry time.
func DecryptSessionToken(token string, key string) (string, int64, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", 0, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", 0, err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", 0, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", 0, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", 0, err
	}

	parts := strings.Split(string(plaintext), "|")
	if len(parts) != 2 {
		return "", 0, errors.New("invalid token format")
	}

	personID := parts[0]
	expiry, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", 0, err
	}

	return personID, expiry, nil
}
