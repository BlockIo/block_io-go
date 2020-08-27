package block_io_go

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/piotrnar/gocoin/lib/btc"
	"golang.org/x/crypto/pbkdf2"
	"log"
)

func SignInputs(ecKey *ECKey, DataToSign string) (string, error) {
	messageHash, _ := hex.DecodeString(DataToSign)
	signature, err := ecKey.Sign(messageHash)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(signature), nil
}

func ParseResponseData(res *resty.Response) (SignatureData, error){
	var withdrawRes SignatureRes
	marshalErr := json.Unmarshal([]byte(res.String()), &withdrawRes)
	if marshalErr != nil {
		return SignatureData{}, marshalErr
	}
	return withdrawRes.Data, nil
}

func ExtractKeyFromEncryptedPassphrase(encryptedData string, b64Key string) (*ECKey, error) {
	aesKey, b64keyErr := base64.StdEncoding.DecodeString(b64Key)
	if (b64keyErr != nil) {
		return nil, b64keyErr
	}

	cipherText, b64CtErr := base64.StdEncoding.DecodeString(encryptedData)
	if (b64CtErr != nil) {
		return nil, b64CtErr
	}

	clearText, decryptErr := AESDecrypt(cipherText, aesKey)
	if (decryptErr != nil) {
		return nil, decryptErr
	}

	seed, hexSeedErr := hex.DecodeString(string(clearText))
	if hexSeedErr != nil {
		return nil, hexSeedErr
	}

	privKey := sha256.Sum256(seed)
	ecKey := NewECKey(privKey, true)
	return ecKey, nil
}

//DEPRECIATED
func ExtractPubKeyFromEncryptedPassphrase(EncryptedData string, B64Key string) string {
	Decrypted := Decrypt(EncryptedData,B64Key)
	Unhexlified, err := hex.DecodeString(string(Decrypted))

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
