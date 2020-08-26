package block_io_go

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/piotrnar/gocoin/lib/btc"
	"golang.org/x/exp/utf8string"
	"log"
)

func FromWIF(privKey string) (string, bool) {
	var extendedKeyBytes []byte = btc.Decodeb58(privKey)
	extendedKeyBytes = extendedKeyBytes[0:34]

	extendedKeyBytes = extendedKeyBytes[1:]
	compressed := false
	if len(extendedKeyBytes) == 33 {
		if extendedKeyBytes[32] != 0x01 {
			log.Fatal(errors.New("Invalid compression flag" + " PrivKey"))
		}
		extendedKeyBytes = extendedKeyBytes[0 : len(extendedKeyBytes)-1]
		compressed = true
	}

	if len(extendedKeyBytes) != 32 {
		log.Fatal(errors.New("Invalid WIF payload length"))
	}

	return hex.EncodeToString(extendedKeyBytes), compressed

}

func PubKeyFromWIF(privKey string) string {

	privKeyFromWifHex, compressed := FromWIF(privKey)
	privKeyFromWifBytes, _ := hex.DecodeString(privKeyFromWifHex)
	result := btc.PublicFromPrivate(privKeyFromWifBytes, compressed)

	return hex.EncodeToString(result)
}

func ExtractKeyFromPassphrase(HexPass string) string {
	Unhexlified, err := hex.DecodeString(HexPass)

	if err != nil {
		log.Fatal(errors.New("Unhexlified Error"))
	}

	Hashed := sha256.Sum256(Unhexlified)
	UsableHashed := Hashed[:]
	return hex.EncodeToString(UsableHashed)
}

func ExtractKeyFromPassphraseString(pass string) string {
	password := []byte(utf8string.NewString(pass).String())
	hashed := sha256.Sum256(password)
	usableHashed := hashed[:]
	return hex.EncodeToString(usableHashed)
}

func ExtractPubKeyFromPassphraseString(pass string) string {
	privKey, _ := hex.DecodeString(ExtractKeyFromPassphraseString(pass))
	result := btc.PublicFromPrivate(privKey, true)
	return hex.EncodeToString(result)
}

func ExtractPubKeyFromPassphrase(HexPass string) string {
	Unhexlified, err := hex.DecodeString(HexPass)

	if err != nil {
		log.Fatal(errors.New("Unhexlified Error"))
	}

	Hashed := sha256.Sum256(Unhexlified)
	UsableHashed := Hashed[:]

	result := btc.PublicFromPrivate(UsableHashed, true)

	return hex.EncodeToString(result)
}
