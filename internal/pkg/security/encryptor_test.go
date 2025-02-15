package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptor(t *testing.T) {
	secretKey := "N1PCdw3M2B1TfJhoaY2mL736p2vCUc47"

	enc := NewEncryptor(secretKey)

	sampleText := "test123"

	encryptedText, err := enc.Encrypt(sampleText)
	assert.NoError(t, err)
	assert.NotNil(t, encryptedText)

	dec, err := enc.Decrypt(*encryptedText)
	assert.NoError(t, err)
	assert.NotNil(t, dec)
	assert.Equal(t, sampleText, *dec)
}
