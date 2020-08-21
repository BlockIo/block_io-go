package sign_request

import (
	"encoding/json"
	"github.com/BlockIo/block_io-go/lib"
)

func Withdraw(pin string, withdrawReqString []byte) lib.SignatureData{

	var signatureRes lib.SignatureData
	err := json.Unmarshal(withdrawReqString, &signatureRes)
	if err != nil {
		panic(err)
	}

	if (signatureRes.ReferenceID == "" || signatureRes.EncryptedPassphrase == lib.EncryptedPassphrase{} ||
		signatureRes.EncryptedPassphrase.Passphrase == "") {
		panic("invalid withdrawal response")
	}
	var encryptedPassphrase = signatureRes.EncryptedPassphrase.Passphrase
	aesKey := lib.PinToAesKey(pin)

	privKey := lib.ExtractKeyFromEncryptedPassphrase(encryptedPassphrase,aesKey)
	pubKey := lib.ExtractPubKeyFromEncryptedPassphrase(encryptedPassphrase,aesKey)

	if pubKey != signatureRes.EncryptedPassphrase.SignerPublicKey {
		panic("Public key mismatch. Invalid Secret PIN detected.")
	}

	for i := 0; i < len(signatureRes.Inputs); i++ {
		for j := 0; j < len(signatureRes.Inputs[i].Signers); j++ {
			signatureRes.Inputs[i].Signers[j].SignedData = lib.SignInputs(privKey, signatureRes.Inputs[i].DataToSign)
		}
	}
	signatureRes.EncryptedPassphrase = lib.EncryptedPassphrase{}

	return signatureRes
}
