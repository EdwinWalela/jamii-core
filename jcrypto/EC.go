package jcrypto

import (
	"crypto/rand"

	"github.com/katzenpost/core/crypto/eddsa"
)

func GenKeyPair(kp *KeyPair, seed string) error {
	var privKey *eddsa.PrivateKey
	var err error

	privKey, err = eddsa.NewKeypair(rand.Reader)

	kp.PrivKey = *privKey
	kp.PubKey = *privKey.PublicKey()

	if err != nil {
		return err
	}

	return nil
}
