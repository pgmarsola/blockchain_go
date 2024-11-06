package utils

import (
	"crypto/ecdsa"
	"encoding/json"
)

const (
	version = byte(0x00)
)

type GenericWallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func (gw GenericWallet) MarshalJSON() ([]byte, error) {
	mapStringAny := map[string]any{
		"PrivateKey": map[string]any{
			"D": gw.PrivateKey.D,
			"PublicKey": map[string]any{
				"X": gw.PrivateKey.PublicKey.X,
				"Y": gw.PrivateKey.PublicKey.Y,
			},
			"X": gw.PrivateKey.X,
			"Y": gw.PrivateKey.Y,
		},
		"PublicKey": gw.PublicKey,
	}
	return json.Marshal(mapStringAny)
}

func MakeGenericWallet() *GenericWallet {
	private, public := NewKeyPair()
	genericWallet := GenericWallet{private, public}

	return &genericWallet
}

func (gw GenericWallet) Address() []byte {
	pubHash := PublicKeyHash(gw.PublicKey)

	versionedHash := append([]byte{version}, pubHash...)
	checkSum := CheckSum(versionedHash)
	fullHash := append(versionedHash, checkSum...)
	address := Base58Encode(fullHash)

	return address
}
