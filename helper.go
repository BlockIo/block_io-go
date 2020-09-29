package block_io_go

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"golang.org/x/crypto/pbkdf2"
	"errors"
)

func ParseResponse(res string) (*BaseResponse, error){
	var resStruct BaseResponse
	marshalErr := json.Unmarshal([]byte(res), &resStruct)

	if marshalErr != nil {
		return nil, marshalErr
	}

	return &resStruct, nil
}

func ParseSignatureResponse(res string) (*SignatureData, error){
  resErr := ValidateResponseJson(res)
	if (resErr != nil) {
		return nil, resErr
	}

	var withdrawRes SignatureRes
	marshalErr := json.Unmarshal([]byte(res), &withdrawRes)
	if marshalErr != nil {
		return nil, marshalErr
	}

	return &withdrawRes.Data, nil
}

func ParseErrorResponse(res string) (*ErrorResponse, error){
	var resStruct ErrorResponse
	marshalErr := json.Unmarshal([]byte(res), &resStruct)
	if marshalErr != nil {
		return nil, marshalErr
	}
	return &resStruct, nil
}

func ValidateResponseJson(responseJson string) error {
	res, baseErr := ParseResponse(responseJson)

	if (baseErr != nil) {
		return baseErr
	}

	if (res.Status != "success") {
		errRes, errErr := ParseErrorResponse(responseJson)
		if (errErr != nil) {
			return errors.New("Cannot parse response from block.io API")
		}
		return errors.New("API ERROR: " + errRes.Data.ErrorMessage)
	}

	return nil
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

	return DeriveKeyFromHex(string(clearText))
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
