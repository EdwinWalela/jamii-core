package jcrypto

import (
	"errors"

	"github.com/katzenpost/core/crypto/eddsa"
)

type KeyPair struct {
	PrivKey eddsa.PrivateKey
	PubKey  eddsa.PublicKey
}

func (kp *KeyPair) Sign(src string) ([]byte, error) {

	if &kp.PrivKey == nil {
		return nil, errors.New("Keypair not generated")
	}

	sig := kp.PrivKey.Sign([]byte(src))

	return sig, nil
}

func (kp *KeyPair) FromBytes(bytes []byte) {
	kp.PrivKey.FromBytes(bytes)
	kp.PubKey = *kp.PrivKey.PublicKey()
}
