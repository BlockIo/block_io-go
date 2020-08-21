package examples

import (
	"encoding/json"
	"fmt"
	"github.com/BlockIo/block_io-go/sign_request"
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

	apiKey := os.Getenv("API_KEY")
	pin := os.Getenv("PIN")

	restClient := resty.New()
	rawWithdrawResponse, err := restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
		"to_address":os.Getenv("TO_ADDRESS"),
		"amount": "0.001",
	}).Post("https://block.io/api/v2/withdraw?api_key=" + apiKey)

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

	signatureRes := sign_request.Withdraw(pin, withdrawResDataString)

	signAndFinalizeReq, pojoErr := json.Marshal(signatureRes)
	if pojoErr != nil {
		panic(pojoErr)
	}

	signAndFinalizeRes, err := restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
		"signature_data": string(signAndFinalizeReq),
	}).
		Post("https://block.io/api/v2/sign_and_finalize_withdrawal?api_key=" + apiKey)

	fmt.Println(signAndFinalizeRes)
}