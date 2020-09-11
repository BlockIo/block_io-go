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
	toAddr := os.Getenv("TO_ADDRESS")
	strWif := os.Getenv("WIF_TO_SWEEP")

	if (apiKey == "" || toAddr == "" || strWif == "") {
		log.Fatal("Sweep requires environment variables API_KEY, TO_ADDRESS and WIF_TO_SWEEP")
	}

	// parse the WIF into an ECKey
	ecKey, wifErr := blockio.FromWIF(strWif)
	if wifErr != nil {
		log.Fatal(wifErr)
	}

	//instantiate a REST client
	restClient := resty.New()

	// post the sweep request to the REST API
	sweepData, sweepErr := restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
			"to_address":   toAddr,
			"public_key":   ecKey.PublicKeyHex(),
		}).Post("https://block.io/api/v2/sweep_from_address?api_key=" + apiKey)

	if sweepErr != nil {
		log.Fatal(sweepErr)
	}

	// sign the request with the ECKey
	signatureReq, signErr := blockio.SignSweepRequest(ecKey, sweepData.String())

	if signErr != nil {
		log.Fatal(signErr)
	}

	// post the resulting signed request to the REST API to broadcast the transaction
	result, postErr := restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
			"signature_data": signatureReq,
		}).Post("https://block.io/api/v2/sign_and_finalize_sweep?api_key=" + apiKey)

	if postErr != nil {
		log.Fatal(postErr)
	}

	log.Println(result)
}
