package auth

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSessionTokenEncryptionDecryption(t *testing.T) {
	os.Setenv("SESSION_KEY", "12345678901234567890123456789012")

	personID := "randomPersonID12345"
	token, err := EncryptSessionToken(personID)
	assert.NoError(t, err)

	decryptedPersonID, expiry, err := DecryptSessionToken(token)
	assert.NoError(t, err)
	assert.Equal(t, personID, decryptedPersonID)

	expiryTime := time.Unix(expiry, 0)
	assert.WithinDuration(t, time.Now().Add(24*time.Hour), expiryTime, 5*time.Second)
}
