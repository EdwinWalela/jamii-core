package jcrypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
)

func GenKeyPair(kp *KeyPair, seed string) error {
	var privKey *ecdsa.PrivateKey
	var err error

	if seed != "" {
		k := new(big.Int)
		k.SetString(seed, 16)

		privKey = new(ecdsa.PrivateKey)
		privKey.PublicKey.Curve = elliptic.P256()
		privKey.D = k

		privKey.PublicKey.X, privKey.PublicKey.Y = elliptic.P256().ScalarBaseMult(k.Bytes())

	} else {
		privKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return err
		}
		kp.PrivKey = *privKey
		kp.PubKey = *&privKey.PublicKey
	}

	return nil
}
