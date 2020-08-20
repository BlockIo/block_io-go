package examples

import (
	"encoding/json"
	"fmt"
	"github.com/BlockIo/block_io-go/lib"
	"github.com/BlockIo/block_io-go/lib/BlockIo"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func WithdrawExample() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	api_key := os.Getenv("API_KEY")
	pin := os.Getenv("PIN")

	restClient := resty.New()
	rawWithdrawResponse, err := restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
		"to_address":os.Getenv("TO_ADDRESS"),
		"amount": "0.001",
	}).Post("https://block.io/api/v2/withdraw?api_key=" + api_key)

	if err != nil {
		panic(err)
	}
	var withdrawResMap map[string]interface{}
	marshalErr := json.Unmarshal([]byte(rawWithdrawResponse.String()), &withdrawResMap)
	if marshalErr != nil {
		panic(marshalErr)
	}
	withdrawResDataString, dataErr := json.Marshal(withdrawResMap["data"])
	if dataErr != nil {
		panic(dataErr)
	}
	var SignatureRes BlockIo.SignatureData
	err = json.Unmarshal(withdrawResDataString, &SignatureRes)
	if err != nil {
		panic(err)
	}

	if (SignatureRes.ReferenceID == "" || SignatureRes.EncryptedPassphrase == BlockIo.EncryptedPassphrase{} ||
		SignatureRes.EncryptedPassphrase.Passphrase == "") {
		panic("invalid withdrawal response")
	}
	var encryptedPassphrase = SignatureRes.EncryptedPassphrase.Passphrase
	 aesKey := lib.PinToAesKey(pin)

	privKey := lib.ExtractKeyFromEncryptedPassphrase(encryptedPassphrase,aesKey)
	pubKey := lib.ExtractPubKeyFromEncryptedPassphrase(encryptedPassphrase,aesKey)

	if pubKey != SignatureRes.EncryptedPassphrase.SignerPublicKey {
		panic("Public key mismatch. Invalid Secret PIN detected.")
	}

	for i := 0; i < len(SignatureRes.Inputs); i++ {
		for j := 0; j < len(SignatureRes.Inputs[i].Signers); j++ {
			SignatureRes.Inputs[i].Signers[j].SignedData = lib.SignInputs(privKey,SignatureRes.Inputs[i].DataToSign)
		}
	}
	SignatureRes.EncryptedPassphrase = BlockIo.EncryptedPassphrase{}
	signAndFinalizeReq, pojoErr := json.Marshal(SignatureRes)
	if pojoErr != nil {
		panic(pojoErr)
	}

	signAndFinalizeRes, err := restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
		"signature_data": string(signAndFinalizeReq),
	}).
		Post("https://block.io/api/v2/sign_and_finalize_withdrawal?api_key=" + api_key)

	fmt.Println(signAndFinalizeRes)
}