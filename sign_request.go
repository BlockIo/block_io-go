package block_io_go

import (
	"encoding/json"
	"encoding/hex"
	"errors"
)

func SignInputs(ecKey *ECKey, DataToSign string) (string, error) {
	messageHash, err := hex.DecodeString(DataToSign)

	if (err != nil) {
		return "", err
	}

	return ecKey.SignHex(messageHash)
}

func signRequest(ecKey *ECKey, reqData *SignatureData) error {
  pubKey := ecKey.PublicKeyHex()

	for i := 0; i < len(reqData.Inputs); i++ {
		for j := 0; j < len(reqData.Inputs[i].Signers); j++ {

			var err error = nil
			if (reqData.Inputs[i].Signers[j].SignerPublicKey != pubKey) {
				continue
			}

			reqData.Inputs[i].Signers[j].SignedData, err = SignInputs(ecKey, reqData.Inputs[i].DataToSign)
			if (err != nil) {
				return err
			}

		}
	}
	reqData.EncryptedPassphrase = EncryptedPassphrase{}

	return nil
}

func SignWithdrawRequest(pin string, withdrawData *SignatureData) (string, error) {

	if (withdrawData.ReferenceID == "" || withdrawData.EncryptedPassphrase == EncryptedPassphrase{} ||
		withdrawData.EncryptedPassphrase.Passphrase == "") {
		return "", errors.New("invalid withdraw response")
	}

	aesKey := PinToAesKey(pin)

	var encryptedPassphrase = withdrawData.EncryptedPassphrase.Passphrase
	ecKey, err := ExtractKeyFromEncryptedPassphrase(encryptedPassphrase, aesKey)
	if (err != nil) {
		return "", err
	}

	pubKeyHex := ecKey.PublicKeyHex()

	if pubKeyHex != withdrawData.EncryptedPassphrase.SignerPublicKey {
		return "", errors.New("Public key mismatch")
	}

	signErr := signRequest(ecKey, withdrawData)
	if signErr != nil {
		return "", signErr
	}

	signAndFinalizeReq, err := json.Marshal(withdrawData)
	if err != nil {
		return "", err
	}

	return string(signAndFinalizeReq), nil
}

// Convenience withdrawal request function that takes a JSON string
func SignWithdrawRequestJson(pin string, withdrawData string) (string, error) {
	withdrawObj, err := ParseSignatureResponse(withdrawData)

	if err != nil {
		return "", err
	}

	return SignWithdrawRequest(pin, withdrawObj)
}

func SignRequestWithKey(eckey *ECKey, sweepReqData *SignatureData) (string, error) {
	if sweepReqData.ReferenceID == "" {
		return "", errors.New("invalid sweep response")
	}

	signErr := signRequest(eckey, sweepReqData)
	if signErr != nil {
		return "", signErr
	}

	signAndFinalizeReq, err := json.Marshal(sweepReqData)
	if err != nil {
		return "", err
	}

	return string(signAndFinalizeReq), nil
}

func SignRequestJson(ecKey *ECKey, sweepData string) (string, error) {
	sweepObj, err := ParseSignatureResponse(sweepData)

	if (err != nil) {
		return "", err
	}

	return SignRequestWithKey(ecKey, sweepObj)
}

func SignRequestWithKeys(ecKeys []*ECKey, dtrustReqData *SignatureData) (string, error) {

	for i := 0; i < len(ecKeys); i++ {
		signErr := signRequest(ecKeys[i], dtrustReqData)
		if (signErr != nil) {
			return "", nil
		}
	}

	output, err := json.Marshal(dtrustReqData)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func SignRequestJsonWithKeys(ecKeys []*ECKey, data string) (string, error) {
	sigObj, err := ParseSignatureResponse(data)

	if (err != nil) {
		return "", err
	}

	return SignRequestWithKeys(ecKeys, sigObj)
}

func SignDtrustRequestWithKeys(ecKeys []*ECKey, data string)(string, error){
	return SignRequestJsonWithKeys(ecKeys, data)
}

func SignDtrustRequestWithKey(ecKey *ECKey, data string)(string, error){
	return SignRequestJson(ecKey, data)
}

func SignSweepRequest(ecKey *ECKey, data string)(string, error){
	return SignRequestJson(ecKey, data)
}
