package jcrypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
)

type KeyPair struct {
	PrivKey ecdsa.PrivateKey
	PubKey  ecdsa.PublicKey
}

func (kp *KeyPair) Sign(src string) ([]byte, error) {

	if &kp.PrivKey == nil {
		return nil, errors.New("Keypair not generated")
	}

	sig, err := ecdsa.SignASN1(rand.Reader, &kp.PrivKey, []byte(src))

	if err != nil {
		return nil, err
	}

	return sig, nil
}

func (kp *KeyPair) Verify(hash, sig []byte) bool {
	return ecdsa.VerifyASN1(&kp.PubKey, hash, sig)
}
