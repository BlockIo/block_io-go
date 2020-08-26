package block_io_go

import (
	"encoding/json"
	"errors"
)

func signRequest(key string, reqData SignatureData) SignatureData {

	for i := 0; i < len(reqData.Inputs); i++ {
		for j := 0; j < len(reqData.Inputs[i].Signers); j++ {
			reqData.Inputs[i].Signers[j].SignedData = SignInputs(key, reqData.Inputs[i].DataToSign)
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
	privKey := ExtractKeyFromEncryptedPassphrase(encryptedPassphrase, aesKey)
	pubKey := ExtractPubKeyFromEncryptedPassphrase(encryptedPassphrase, aesKey)

	if pubKey != withdrawData.EncryptedPassphrase.SignerPublicKey {
		return "", errors.New("Public key mismatch. Invalid Secret PIN detected.")
	}

	signedWithdrawReqData := signRequest(privKey, withdrawData)

	signAndFinalizeReq, err := json.Marshal(signedWithdrawReqData)
	if err != nil {
		return "", err
	}
	return string(signAndFinalizeReq), nil
}

func SignSweepRequest(privKey string, sweepReqData SignatureData) (string, error) {
	if sweepReqData.ReferenceID == "" {
		return "", errors.New("invalid sweep response")
	}
	keyFromWif, _ := FromWIF(privKey)
	signedSweepReqData := signRequest(keyFromWif, sweepReqData)
	signAndFinalizeReq, err := json.Marshal(signedSweepReqData)
	if err != nil {
		return "", err
	}
	return string(signAndFinalizeReq), nil
}
