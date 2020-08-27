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
