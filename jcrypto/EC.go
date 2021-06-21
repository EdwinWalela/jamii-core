package jcrypto

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"

	"github.com/katzenpost/core/crypto/eddsa"
)

func GenKeyPair(kp *KeyPair, secret string) error {
	var privKey *eddsa.PrivateKey
	var err error
	reader := rand.Reader

	privKey, err = eddsa.NewKeypair(reader)

	kp.PrivKey = *privKey
	kp.PubKey = *privKey.PublicKey()

	writeErr := ioutil.WriteFile("priv.dat", kp.PrivKey.Bytes(), 0644)

	fmt.Println(writeErr)

	if err != nil {
		return err
	}

	return nil
}

func ReadKeyPair(kp *KeyPair, path string) {
	dat, _ := ioutil.ReadFile("priv.dat")
	kp.FromBytes(dat)
}

func PubKeyFromBytes(key []byte, kp *KeyPair) {
	kp.PubKey.FromBytes(key)
}

func VerifySig(sig, msg []byte, address eddsa.PublicKey) bool {
	return address.Verify(sig, msg)
}
