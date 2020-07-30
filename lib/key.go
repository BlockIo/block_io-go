package lib

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/piotrnar/gocoin/lib/btc"
	"golang.org/x/exp/utf8string"
	"log"
)

func FromWIF(PrivKey string) string {
	var ExtendedKeyBytes []byte = btc.Decodeb58(PrivKey)
	ExtendedKeyBytes = ExtendedKeyBytes[0:34]

	ExtendedKeyBytes = ExtendedKeyBytes[1:]

	if len(ExtendedKeyBytes) == 33 {
		if ExtendedKeyBytes[32] != 0x01 {
			log.Fatal(errors.New("Invalid compression flag" + " PrivKey"))
		}
		ExtendedKeyBytes = ExtendedKeyBytes[0 : len(ExtendedKeyBytes)-1]
	}

	if len(ExtendedKeyBytes) != 32 {
		log.Fatal(errors.New("Invalid WIF payload length"))
	}

	return hex.EncodeToString(ExtendedKeyBytes)

}

func PubKeyFromWIF(PrivKey string) string {

	var ExtendedKeyBytes []byte = btc.Decodeb58(PrivKey)
	ExtendedKeyBytes = ExtendedKeyBytes[0:34]
	compressed := false
	ExtendedKeyBytes = ExtendedKeyBytes[1:]

	if len(ExtendedKeyBytes) == 33 {
		if ExtendedKeyBytes[32] != 0x01 {
			log.Fatal(errors.New("Invalid compression flag" + " PrivKey"))
		}
		ExtendedKeyBytes = ExtendedKeyBytes[0 : len(ExtendedKeyBytes)-1]
		compressed = true
	}

	if len(ExtendedKeyBytes) != 32 {
		log.Fatal(errors.New("Invalid WIF payload length"))
	}

	result := btc.PublicFromPrivate(ExtendedKeyBytes, compressed)

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
	UsableHashed := hashed[:]
	return hex.EncodeToString(UsableHashed)
}

func ExtractPubKeyFromPassphraseString(pass string) string {
	password := []byte(utf8string.NewString(pass).String())
	hashed := sha256.Sum256(password)
	UsableHashed := hashed[:]
	result := btc.PublicFromPrivate(UsableHashed, true)
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
