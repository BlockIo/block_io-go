package block_io_go

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/piotrnar/gocoin/lib/btc"
	"github.com/btcsuite/btcd/btcec"
	"log"
)

type ECKey struct {
	d [32]byte
	Compressed bool
}

func NewECKey (d [32]byte, compressed bool) *ECKey {
	return &ECKey{
		d: d,
		Compressed: compressed,
	}
}

func (eck *ECKey) PrivateKey() []byte {
	return eck.d[:]
}

func (eck *ECKey) PrivateKeyHex() string {
	return hex.EncodeToString(eck.PrivateKey())
}

func (eck *ECKey) PublicKey() []byte {
	return btc.PublicFromPrivate(eck.PrivateKey(), eck.Compressed)
}

func (eck *ECKey) PublicKeyHex() string {
	return hex.EncodeToString(eck.PublicKey())
}

func (eck *ECKey) Sign(message []byte) ([]byte, error) {
	//TODO use gocoin instead of btcsuite/btcd
	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), eck.d[:])
	signature, err := privKey.Sign(message)
	if (err != nil) {
		return nil, err
	}

	return signature.Serialize(), nil
}

func (eck *ECKey) SignHex(message []byte) (string, error) {
	signature, err := eck.Sign(message)
	if (err != nil) {
		return "", err
	}
	return hex.EncodeToString(signature), nil
}

func FromWIF(privKey string) (*ECKey, error) {

	var keyBytes [32]byte
	var extendedKeyBytes []byte = btc.Decodeb58(privKey)
	extendedKeyBytes = extendedKeyBytes[0:34]
	extendedKeyBytes = extendedKeyBytes[1:]

	compressed := false
	if len(extendedKeyBytes) == 33 {
		if extendedKeyBytes[32] != 0x01 {
			return nil, errors.New("ECKey/FromWIF: Invalid compression flag")
		}
		compressed = true
	} else if len(extendedKeyBytes) != 32 {
		return nil, errors.New("Invalid WIF payload")
	}

	copy(keyBytes[:], extendedKeyBytes)
	eckey := NewECKey(keyBytes, compressed)

	return eckey, nil
}

func ExtractKeyFromPassphrase(HexPass string) *ECKey {
	Unhexlified, err := hex.DecodeString(HexPass)

	if err != nil {
		log.Fatal(errors.New("Unhexlified Error"))
	}

	hashed := sha256.Sum256(Unhexlified)
	return NewECKey(hashed, true)
}

func ExtractKeyFromPassphraseString(pass string) *ECKey {
	password := []byte(pass)
	hashed := sha256.Sum256(password)
	return NewECKey(hashed, true)
}
