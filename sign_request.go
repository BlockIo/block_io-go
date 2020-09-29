package block_io_go

import (
	"encoding/json"
	"encoding/hex"
	"errors"
)

// Sign any hex string
func SignInputs(ecKey *ECKey, hexData string) (string, error) {
	data, err := hex.DecodeString(hexData)

	if (err != nil) {
		return "", err
	}

	return ecKey.SignHex(data)
}

func signRequest(ecKey *ECKey, reqData *SignatureData) error {
  pubKey := ecKey.PublicKeyHex()

	for i := 0; i < len(reqData.Inputs); i++ {
		for j := 0; j < len(reqData.Inputs[i].Signers); j++ {

			if (reqData.Inputs[i].Signers[j].SignerPublicKey != pubKey) {
				continue
			}

			var err error = nil
			reqData.Inputs[i].Signers[j].SignedData, err = SignInputs(ecKey, reqData.Inputs[i].DataToSign)

			if (err != nil) {
				return err
			}

		}
	}

	return nil
}

// Sign a withdrawal request with a pin
func SignWithdrawRequest(pin string, withdrawData *SignatureData) (string, error) {

	if (withdrawData.EncryptedPassphrase == EncryptedPassphrase{} ||
		withdrawData.EncryptedPassphrase.Passphrase == "") {
		return "", errors.New("Withdrawal sign request is missing encrypted passphrase")
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

	withdrawData.EncryptedPassphrase = EncryptedPassphrase{}

	return SignRequestWithKey(ecKey, withdrawData)
}

// Convenience withdrawal request function that takes a JSON string
func SignWithdrawRequestJson(pin string, withdrawData string) (string, error) {
	withdrawObj, err := ParseSignatureResponse(withdrawData)

	if err != nil {
		return "", err
	}

	return SignWithdrawRequest(pin, withdrawObj)
}

// Sign a withdrawal request with a custom ECKey
func SignRequestWithKey(eckey *ECKey, sigRequest *SignatureData) (string, error) {
	if sigRequest.ReferenceID == "" {
		return "", errors.New("Signing request is missing referenceId")
	}

	signErr := signRequest(eckey, sigRequest)
	if signErr != nil {
		return "", signErr
	}

	jsonData, err := json.Marshal(sigRequest)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// Sign a JSON string withdrawal request with a custom ECKey
func SignRequestJsonWithKey(ecKey *ECKey, data string) (string, error) {
	sigRequest, err := ParseSignatureResponse(data)

	if (err != nil) {
		return "", err
	}

	return SignRequestWithKey(ecKey, sigRequest)
}

// Sign a withdrawal request with a set of custom ECKeys
func SignRequestWithKeys(ecKeys []*ECKey, sigRequest *SignatureData) (string, error) {

	if sigRequest.ReferenceID == "" {
		return "", errors.New("Signing request is missing referenceId")
	}

	for i := 0; i < len(ecKeys); i++ {
		signErr := signRequest(ecKeys[i], sigRequest)
		if (signErr != nil) {
			return "", nil
		}
	}

	output, err := json.Marshal(sigRequest)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// Sign a JSON string withdrawal request with a set of custom ECKeys
func SignRequestJsonWithKeys(ecKeys []*ECKey, data string) (string, error) {
	sigObj, err := ParseSignatureResponse(data)

	if (err != nil) {
		return "", err
	}

	return SignRequestWithKeys(ecKeys, sigObj)
}

//DEPRECIATED. Use SignRequestWithKeys or SignRequestJsonWithKeys
func SignDtrustRequest(ecKeys []*ECKey, dtrustReqData SignatureData) (string, error) {
	return SignRequestWithKeys(ecKeys, &dtrustReqData)
}

//DEPRECIATED. Use SignRequestWithKey or SignRequestJsonWithKey
func SignSweepRequest(ecKey *ECKey, sweepReqData SignatureData) (string, error) {
	return SignRequestWithKey(ecKey, &sweepReqData)
}
