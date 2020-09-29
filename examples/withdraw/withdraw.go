package main

import (
	blockio "github.com/BlockIo/block_io-go"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	// load vars from .env file if it's there
	godotenv.Load(".env")

	// load environment vars
	apiKey := os.Getenv("API_KEY")
	pin := os.Getenv("PIN")
	destinationAddress := os.Getenv("TO_ADDRESS")

	// instantiate a REST client
	restClient := resty.New()

	// post the withdraw request to the REST API
	withdrawRequest, withdrawErr := restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
		"to_address": destinationAddress,
		"amount": "0.01",
	}).Post("https://block.io/api/v2/withdraw?api_key=" + apiKey)

	if withdrawErr != nil {
		log.Fatal(withdrawErr)
	}

	// sign the request
	signatureReq, signErr := blockio.SignWithdrawRequestJson(pin, withdrawRequest.String())

	if signErr != nil {
		log.Fatal(signErr)
	}

	// post the resulting signed request to the REST API to broadcast the transaction
	result, postErr := restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
		"signature_data": signatureReq,
	}).Post("https://block.io/api/v2/sign_and_finalize_withdrawal?api_key=" + apiKey)

	if postErr != nil {
		log.Fatal(postErr)
	}

	log.Println(result)
}
