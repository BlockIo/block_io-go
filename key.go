package block_io_go

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"log"
)

type ECKey struct {
	d *btcec.PrivateKey
	Compressed bool
}

func NewECKey (d [32]byte, compressed bool) *ECKey {
	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), d[:])
	return &ECKey{
		d: privKey,
		Compressed: compressed,
	}
}

func (eck *ECKey) PrivateKey() []byte {
	return eck.d.Serialize()
}

func (eck *ECKey) PrivateKeyHex() string {
	return hex.EncodeToString(eck.PrivateKey())
}

func (eck *ECKey) PublicKey() []byte {
	if (eck.Compressed) {
		return eck.d.PubKey().SerializeCompressed()
	}
	return eck.d.PubKey().SerializeUncompressed()
}

func (eck *ECKey) PublicKeyHex() string {
	return hex.EncodeToString(eck.PublicKey())
}

func (eck *ECKey) Sign(message []byte) ([]byte, error) {
	signature, err := eck.d.Sign(message)
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

func FromWIF(strWif string) (*ECKey, error) {
	wif, err := btcutil.DecodeWIF(strWif)
	if (err != nil) {
		return nil, err
	}

	eckey := &ECKey{wif.PrivKey, wif.CompressPubKey}
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
