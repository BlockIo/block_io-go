package main

import (
	"fmt"
	blockio "github.com/BlockIo/block_io-go"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	godotenv.Load(".env")

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
		log.Fatal(err)
	}

	fmt.Println("Raw withdraw response: ")
	fmt.Println(rawWithdrawResponse)

	withdrawData, withdrawDataErr := blockio.ParseResponseData(rawWithdrawResponse.String())

	if withdrawDataErr != nil {
		log.Fatal(withdrawDataErr)
	}

	signatureReq, signWithdrawReqErr := blockio.SignWithdrawRequest(pin, withdrawData)

	if signWithdrawReqErr != nil {
		log.Fatal(signWithdrawReqErr)
	}

	signAndFinalizeRes, err := restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
		"signature_data": signatureReq,
	}).Post("https://block.io/api/v2/sign_and_finalize_withdrawal?api_key=" + apiKey)

	fmt.Println(signAndFinalizeRes)
}
