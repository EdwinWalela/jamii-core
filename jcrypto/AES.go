package jcrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

func AESEncrypt(secret []byte, source []byte) ([]byte, error) {

	// s := len(secret)

	// secret = append(secret, )

	c, err := aes.NewCipher(secret)

	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)

	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, source, nil), nil
}
