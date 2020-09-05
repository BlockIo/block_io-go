package block_io_go

import (
	"encoding/json"
	"encoding/hex"
	"errors"
)

func SignInputs(ecKey *ECKey, DataToSign string) (string, error) {
	messageHash, _ := hex.DecodeString(DataToSign)
	return ecKey.SignHex(messageHash)
}

func signRequest(ecKey *ECKey, reqData SignatureData) SignatureData {

	for i := 0; i < len(reqData.Inputs); i++ {
		for j := 0; j < len(reqData.Inputs[i].Signers); j++ {
			//TODO handle my errors
			reqData.Inputs[i].Signers[j].SignedData, _ = SignInputs(ecKey, reqData.Inputs[i].DataToSign)
		}
	}
	reqData.EncryptedPassphrase = EncryptedPassphrase{}

	return reqData
}

func SignWithdrawRequest(pin string, withdrawData SignatureData) (string, error) {

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

	signedWithdrawReqData := signRequest(ecKey, withdrawData)

	signAndFinalizeReq, err := json.Marshal(signedWithdrawReqData)
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

func SignSweepRequest(eckey *ECKey, sweepReqData SignatureData) (string, error) {
	if sweepReqData.ReferenceID == "" {
		return "", errors.New("invalid sweep response")
	}

	signedSweepReqData := signRequest(eckey, sweepReqData)
	signAndFinalizeReq, err := json.Marshal(signedSweepReqData)
	if err != nil {
		return "", err
	}
	return string(signAndFinalizeReq), nil
}

func SignSweepRequestJson(ecKey *ECKey, sweepData string) (string, error) {
	sweepObj, err := ParseSignatureResponse(sweepData)

	if (err != nil) {
		return "", err
	}

	return SignSweepRequest(ecKey, sweepObj)
}

func SignDtrustRequest(ecKeys []*ECKey, dtrustReqData SignatureData) (string, error) {

	for i := 0; i < len(dtrustReqData.Inputs); i++ {
		for j := 0; j < len(dtrustReqData.Inputs[i].Signers); j++ {
			//TODO handle my errors
			dtrustReqData.Inputs[i].Signers[j].SignedData, _ = SignInputs(ecKeys[j], dtrustReqData.Inputs[i].DataToSign)
		}
	}

	signAndFinalizeReq, err := json.Marshal(dtrustReqData)
	if err != nil {
		return "", err
	}
	return string(signAndFinalizeReq), nil
}
