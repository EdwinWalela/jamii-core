package primitives

import (
	"encoding/base64"
	"strconv"
	"strings"

	"github.com/edwinwalela/jamii-core/jcrypto"
	"github.com/katzenpost/core/crypto/eddsa"
)

type Vote struct {
	/**
	Vote represents a tx in the blockchain
	Mobile clients submits vote via TCP socket connection & node packs it into a Vote
	**/

	Address    eddsa.PublicKey   // public key of tx initiator
	Candidates []eddsa.PublicKey // list of candidtes selected by initiator
	Signature  []byte            // signature by the initiator's public key
	Hash       string            // hash of the tx
	Timestamp  uint64            // Unix timestamp of tx in seconds
}

/**
	string data from client ->	pubkey64|sign64|hash|candidate,candidate,candidate|timestamp
	This is broken down by node
**/

func (v *Vote) UnpackClientString(data string) error {

	vote := strings.Split(data, "|")

	var sig64, pub64 string
	var uintErr error

	candidates := make([]eddsa.PublicKey, 1)

	for i, val := range vote {
		switch i {
		case 0: // Extract vote hash
			v.Hash = val
		case 1: // Extract sign64
			sig64 = val
		case 2: // Extract pub64
			pub64 = val
		case 3: // Extract candidates []pub64
			for _, candidatePub64 := range strings.Split(val, ",") {
				decodedPub, pubErr := base64.StdEncoding.DecodeString(candidatePub64)

				if pubErr != nil {
					return pubErr
				}
				kp := &jcrypto.KeyPair{}
				jcrypto.PubKeyFromBytes(decodedPub, kp)
				candidates = append(candidates, kp.PubKey)
			}
		case 4:
			v.Timestamp, uintErr = strconv.ParseUint(val, 10, 64)
			if uintErr != nil {
				return uintErr
			}
		}
	}

	decodedPub, pubErr := base64.StdEncoding.DecodeString(pub64)
	decodedSig, sigErr := base64.StdEncoding.DecodeString(sig64)

	if pubErr != nil {
		return pubErr
	}

	if sigErr != nil {
		return sigErr
	}

	kp := &jcrypto.KeyPair{}

	jcrypto.PubKeyFromBytes(decodedPub, kp)

	v.Address = kp.PubKey
	v.Signature = decodedSig
	v.Candidates = candidates

	return nil
}
