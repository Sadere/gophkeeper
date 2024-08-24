package utils

import "crypto/rand"

func GenerateRandom(size int) ([]byte, error) {
	// generating cryptographically strong random bytes in b
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
