package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

type Encryptor interface {
	Encrypt(plaintext string) (*string, error)
	Decrypt(ciphertext string) (*string, error)
}

type encryptor struct {
	cipherBlock cipher.Block
	aead        cipher.AEAD
}

func NewEncryptor(secretKey string) Encryptor {
	// Load the secret key from a secure location.
	aes, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		panic(err)
	}

	return &encryptor{
		cipherBlock: aes,
		aead:        gcm,
	}
}

func NewEncryptorFromEnv() Encryptor {
	// Load the secret key from a secure location.
	secretKey := viper.GetString("storage.encryption.key")

	aes, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		panic(err)
	}

	return &encryptor{
		cipherBlock: aes,
		aead:        gcm,
	}
}

func (e *encryptor) Encrypt(plaintext string) (*string, error) {
	// We need a 12-byte nonce for GCM (modifiable if you use cipher.NewGCMWithNonceSize())
	// A nonce should always be randomly generated for every encryption.
	nonce := make([]byte, e.aead.NonceSize())
	_, err := rand.Read(nonce)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate nonce")
	}

	// ciphertext here is actually nonce+ciphertext
	// So that when we decrypt, just knowing the nonce size
	// is enough to separate it from the ciphertext.
	ciphertext := e.aead.Seal(nonce, nonce, []byte(plaintext), nil)

	return lo.ToPtr(base64.StdEncoding.EncodeToString(ciphertext)), nil
}

func (e *encryptor) Decrypt(ciphertext string) (*string, error) {
	decoded, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode ciphertext")
	}

	// Since we know the ciphertext is actually nonce+ciphertext
	// And len(nonce) == NonceSize(). We can separate the two.
	nonceSize := e.aead.NonceSize()
	nonce := decoded[:nonceSize]
	actualCiphertext := decoded[nonceSize:]

	plaintext, err := e.aead.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt ciphertext")
	}

	return lo.ToPtr(string(plaintext)), nil
}
