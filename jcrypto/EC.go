package jcrypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"io/ioutil"
)

func GenKeyPair(kp *KeyPair) error {

	reader := rand.Reader

	pubKey, privKey, err := ed25519.GenerateKey(reader)

	kp.PrivKey = privKey
	kp.PubKey = pubKey

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
	// kp.PubKey.FromBytes(key)
	// ed25519.PublicKey
}

func VerifySig(sig, msg []byte, address ed25519.PublicKey) bool {
	return ed25519.Verify(address, msg, sig)

}
