package main

import (
	"encoding/json"
)

func signRequest(key string, reqString []byte, requestType string) SignatureData {

	var signatureRes SignatureData
	err := json.Unmarshal(reqString, &signatureRes)
	if err != nil {
		panic(err)
	}

	if signatureRes.ReferenceID == "" {
		panic("invalid " + requestType + " response")
	}

	switch requestType {
	case "withdraw":
		if (signatureRes.EncryptedPassphrase == EncryptedPassphrase{} ||
			signatureRes.EncryptedPassphrase.Passphrase == "") {
			panic("invalid withdrawal response")
		}
		var encryptedPassphrase = signatureRes.EncryptedPassphrase.Passphrase

		privKey := ExtractKeyFromEncryptedPassphrase(encryptedPassphrase, key)
		pubKey := ExtractPubKeyFromEncryptedPassphrase(encryptedPassphrase, key)

		if pubKey != signatureRes.EncryptedPassphrase.SignerPublicKey {
			panic("Public key mismatch. Invalid Secret PIN detected.")
		}
		key = privKey
	}

	for i := 0; i < len(signatureRes.Inputs); i++ {
		for j := 0; j < len(signatureRes.Inputs[i].Signers); j++ {
			signatureRes.Inputs[i].Signers[j].SignedData = SignInputs(key, signatureRes.Inputs[i].DataToSign)
		}
	}
	signatureRes.EncryptedPassphrase = EncryptedPassphrase{}

	return signatureRes
}

func SignWithdrawRequest(pin string, withdrawReqString []byte) SignatureData{

	aesKey := PinToAesKey(pin)
	return signRequest(aesKey, withdrawReqString, "withdraw")
}

func SignSweepRequest(privKey string, sweepReqString []byte) SignatureData{
	keyFromWif, _ := FromWIF(privKey)
	return signRequest(keyFromWif, sweepReqString, "sweep")
}
