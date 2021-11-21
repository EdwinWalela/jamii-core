package jcrypto

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"log"
)

type KeyPair struct {
	PrivKey ed25519.PrivateKey
	PubKey  ed25519.PublicKey
}

func (kp *KeyPair) Sign(src string) ([]byte, error) {

	if &kp.PrivKey == nil {
		return nil, errors.New("Keypair not generated")
	}
	reader := rand.Reader

	sig, err := kp.PrivKey.Sign(reader, []byte(src), crypto.Hash(0))

	if err != nil {
		log.Println(err)
	}
	return sig, nil
}

func (kp *KeyPair) FromBytes(bytes []byte) {
	// kp.PrivKey.FromBytes(bytes)
	// kp.PubKey = *kp.PrivKey.PublicKey()
}
