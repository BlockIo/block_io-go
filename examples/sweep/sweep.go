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
	toAddr := os.Getenv("TO_ADDRESS")
	strWif := os.Getenv("WIF_TO_SWEEP")

	if (apiKey == "" || toAddr == "" || strWif == "") {
		log.Fatal("Sweep requires environment variables API_KEY, TO_ADDRESS and WIF_TO_SWEEP")
	}

	ecKey, wifErr := blockio.FromWIF(strWif)
	if wifErr != nil {
		log.Fatal(wifErr)
	}

	restClient := resty.New()
	rawSweepResponse, err := restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
			"to_address":   toAddr,
			"public_key":   ecKey.PublicKeyHex(),
		}).Post("https://block.io/api/v2/sweep_from_address?api_key=" + apiKey)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Raw sweep response: ")
	fmt.Println(rawSweepResponse)

	sweepData, sweepDataErr := blockio.ParseSignatureResponse(rawSweepResponse.String())

	if sweepDataErr != nil {
		log.Fatal(sweepDataErr)
	}

	signatureReq, signSweepReqErr := blockio.SignSweepRequest(ecKey, sweepData)

	if signSweepReqErr != nil {
		log.Fatal(signSweepReqErr)
	}

	signAndFinalizeRes, err := restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
			"signature_data": signatureReq,
		}).Post("https://block.io/api/v2/sign_and_finalize_sweep?api_key=" + apiKey)

	fmt.Println(signAndFinalizeRes)
}
