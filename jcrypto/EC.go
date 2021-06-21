package jcrypto

import (
	"crypto/rand"
	"io/ioutil"

	"github.com/katzenpost/core/crypto/eddsa"
)

func GenKeyPair(kp *KeyPair) error {
	var privKey *eddsa.PrivateKey
	var err error
	reader := rand.Reader

	privKey, err = eddsa.NewKeypair(reader)

	kp.PrivKey = *privKey
	kp.PubKey = *privKey.PublicKey()

	if err != nil {
		return err
	}

	return nil
}

func ReadKeyPair(kp *KeyPair, path string) error {
	dat, e := ioutil.ReadFile(path)
	if e != nil {
		return e
	}
	kp.FromBytes(dat)
	return nil
}

func PubKeyFromBytes(key []byte, kp *KeyPair) {
	kp.PubKey.FromBytes(key)
}

func VerifySig(sig, msg []byte, address eddsa.PublicKey) bool {
	return address.Verify(sig, msg)
}
