package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/Sadere/gophkeeper/internal/server/utils"
	"github.com/Sadere/gophkeeper/pkg/constants"

	"golang.org/x/crypto/pbkdf2"
)

// Encrypts plaintext data using AES-GCM with a key derived from master password
func Encrypt(password string, plaintext []byte) ([]byte, error) {
	key, salt := deriveKey(password, nil)

	// creating AES block
	AESBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// creating GCM
	GCM, err := cipher.NewGCM(AESBlock)
	if err != nil {
		return nil, err
	}

	// generating nonce
	nonce, err := utils.GenerateRandom(GCM.NonceSize())
	if err != nil {
		return nil, err
	}

	// encrypt data
	encrypted := GCM.Seal(nonce, nonce, plaintext, nil)

	// store salt alongside encrypted data
	encrypted = append(encrypted, salt...)

	return encrypted, nil
}

// Decrypts data using master password
func Decrypt(password string, encrypted []byte) ([]byte, error) {
	// extract salt
	saltIdx := len(encrypted) - constants.SaltLen
	salt := encrypted[saltIdx:]

	encrypted = encrypted[:saltIdx]

	key, _ := deriveKey(password, salt)

	// creating AES block
	AESBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// creating GCM
	GCM, err := cipher.NewGCM(AESBlock)
	if err != nil {
		return nil, err
	}

	// extract nonce
	nonce := encrypted[:GCM.NonceSize()]
	encrypted = encrypted[GCM.NonceSize():]

	// encrypt data
	decrypted, err := GCM.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return decrypted, nil
}

func deriveKey(password string, salt []byte) ([]byte, []byte) {
	if len(salt) == 0 {
		salt = make([]byte, constants.SaltLen)
		rand.Read(salt)
	}
	return pbkdf2.Key([]byte(password), salt, 4096, 32, sha256.New), salt
}
