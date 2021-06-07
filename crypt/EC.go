package crypt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

func GenKeyPair(kp *KeyPair) error {
	var err error
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	kp.PrivKey = *privKey
	kp.PubKey = *&privKey.PublicKey

	if err != nil {
		return err
	}

	return nil
}
