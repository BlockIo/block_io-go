package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func SweepExample() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	apiKey := os.Getenv("API_KEY")
	toAddr := os.Getenv("TO_ADDRESS")
	privKey := os.Getenv("PRIVATE_KEY_FROM_ADDRESS")
	fromAddr := os.Getenv("FROM_ADDRESS")

	restClient := resty.New()
	rawSweepResponse, err := restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
			"to_address": toAddr,
			"public_key": GetPublicKey(privKey),
			"from_address": fromAddr,
		}).Post("https://block.io/api/v2/sweep_from_address?api_key=" + apiKey)

	if err != nil {
		panic(err)
	}
	var sweepResMap map[string]interface{}
	marshalErr := json.Unmarshal([]byte(rawSweepResponse.String()), &sweepResMap)
	if marshalErr != nil {
		panic(marshalErr)
	}
	sweepResDataString, dataErr := json.Marshal(sweepResMap["data"])
	if dataErr != nil {
		panic(dataErr)
	}

	signatureRes := SignSweepRequest(privKey, sweepResDataString)

	signAndFinalizeReq, pojoErr := json.Marshal(signatureRes)
	if pojoErr != nil {
		panic(pojoErr)
	}

	signAndFinalizeRes, err := restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
			"signature_data": string(signAndFinalizeReq),
		}).
		Post("https://block.io/api/v2/sign_and_finalize_sweep?api_key=" + apiKey)

	fmt.Println(signAndFinalizeRes)
}