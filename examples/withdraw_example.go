package main

import (
	"fmt"
	blockIO "github.com/BlockIo/block_io-go"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
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

	fmt.Println("Raw withdraw response: ")
	fmt.Println(rawWithdrawResponse)

	withdrawData, withdrawDataErr := blockIO.ParseResponseData(rawWithdrawResponse)

	if withdrawDataErr != nil {
		panic(withdrawDataErr)
	}

	signatureReq, signWithdrawReqErr := blockIO.SignWithdrawRequest(pin, withdrawData)

	if signWithdrawReqErr != nil {
		panic(signWithdrawReqErr)
	}

	signAndFinalizeRes, err := restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
		"signature_data": signatureReq,
	}).
		Post("https://block.io/api/v2/sign_and_finalize_withdrawal?api_key=" + apiKey)

	fmt.Println(signAndFinalizeRes)
}