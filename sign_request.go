package blockIO

import "encoding/json"

func signRequest(key string, reqData SignatureData, requestType string) SignatureData {

	if reqData.ReferenceID == "" {
		panic("invalid " + requestType + " response")
	}

	switch requestType {
	case "withdraw":
		if (reqData.EncryptedPassphrase == EncryptedPassphrase{} ||
			reqData.EncryptedPassphrase.Passphrase == "") {
			panic("invalid withdrawal response")
		}
		var encryptedPassphrase = reqData.EncryptedPassphrase.Passphrase

		privKey := ExtractKeyFromEncryptedPassphrase(encryptedPassphrase, key)
		pubKey := ExtractPubKeyFromEncryptedPassphrase(encryptedPassphrase, key)

		if pubKey != reqData.EncryptedPassphrase.SignerPublicKey {
			panic("Public key mismatch. Invalid Secret PIN detected.")
		}
		key = privKey
	}

	for i := 0; i < len(reqData.Inputs); i++ {
		for j := 0; j < len(reqData.Inputs[i].Signers); j++ {
			reqData.Inputs[i].Signers[j].SignedData = SignInputs(key, reqData.Inputs[i].DataToSign)
		}
	}
	reqData.EncryptedPassphrase = EncryptedPassphrase{}

	return reqData
}

func SignWithdrawRequest(pin string, withdrawData SignatureData) (string, error) {

	aesKey := PinToAesKey(pin)
	signedWithdrawReqData := signRequest(aesKey, withdrawData, "withdraw")
	signAndFinalizeReq, err := json.Marshal(signedWithdrawReqData)
	if err != nil {
		return "", err
	}
	return string(signAndFinalizeReq), nil
}

func SignSweepRequest(privKey string, sweepReqData SignatureData) (string, error) {
	keyFromWif, _ := FromWIF(privKey)
	signedSweepReqData := signRequest(keyFromWif, sweepReqData, "sweep")
	signAndFinalizeReq, err := json.Marshal(signedSweepReqData)
	if err != nil {
		return "", err
	}
	return string(signAndFinalizeReq), nil
}
