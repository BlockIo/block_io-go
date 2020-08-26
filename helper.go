package block_io_go

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/go-resty/resty/v2"
	"github.com/piotrnar/gocoin/lib/btc"
	"golang.org/x/crypto/pbkdf2"
	"log"
)

func SignInputs(PrivKey string, DataToSign string) string {
	// Decode a hex-encoded private key.
	pkBytes, err := hex.DecodeString(PrivKey)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)
	pubKey = pubKey
	// Sign a message using the private key.
	message := DataToSign
	messageHash, _ := hex.DecodeString(message)
	signature, err := privKey.Sign(messageHash)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return hex.EncodeToString(signature.Serialize())
}

func ParseResponseData(res *resty.Response) (SignatureData, error){
	var withdrawRes SignatureRes
	marshalErr := json.Unmarshal([]byte(res.String()), &withdrawRes)
	if marshalErr != nil {
		return SignatureData{}, marshalErr
	}
	return withdrawRes.Data, nil
}

func ExtractKeyFromEncryptedPassphrase(EncryptedData string, B64Key string) string {
	Decrypted := Decrypt(EncryptedData,B64Key)
	Unhexlified, err := hex.DecodeString(Decrypted)

	if err != nil {
		log.Fatal(errors.New("Unhexlified Error"))
	}

	Hashed := sha256.Sum256(Unhexlified)
	UsableHashed := Hashed[:]
	return hex.EncodeToString(UsableHashed)
}

func ExtractPubKeyFromEncryptedPassphrase(EncryptedData string, B64Key string) string {
	Decrypted := Decrypt(EncryptedData,B64Key)
	Unhexlified, err := hex.DecodeString(Decrypted)

	if err != nil {
		log.Fatal(errors.New("Unhexlified Error"))
	}

	Hashed := sha256.Sum256(Unhexlified)
	UsableHashed := Hashed[:]

	result := btc.PublicFromPrivate(UsableHashed, true)

	return hex.EncodeToString(result)
}

func PinToAesKey(pin string) string {
	var saltOld []byte = make([]byte, 0)
	var salt [1024]byte;
	for i := 0; i < 1024; i++ {
		salt[i] = 0
	}
	pinBytes := []byte(pin)

	firstHash := hex.EncodeToString(pbkdf2.Key(pinBytes, saltOld, 1024, 16, sha256.New))

	firstHashBytes := []byte(firstHash)

	key := pbkdf2.Key(firstHashBytes, saltOld, 1024, 32, sha256.New)

	return base64.StdEncoding.EncodeToString(key)
}

func SHA256_hash(ba []byte) []byte {
	sha := sha256.Sum256(ba)
	return sha[:]
}
