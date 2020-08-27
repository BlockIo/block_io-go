package main

import (
	"fmt"
	blockio "github.com/BlockIo/block_io-go"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"math/rand"
	"os"
	"strings"
)

func main(){
	godotenv.Load(".env")
	apiKey := os.Getenv("API_KEY")
	pin := os.Getenv("PIN")

	dtrustAddress := ""
	dtrustAddressLabel := "dTrust1" + fmt.Sprintf("%f", rand.ExpFloat64())

	privKeys := []*blockio.ECKey{
		blockio.ExtractKeyFromPassphraseString("verysecretkey1"),
		blockio.ExtractKeyFromPassphraseString("verysecretkey2"),
		blockio.ExtractKeyFromPassphraseString("verysecretkey3"),
		blockio.ExtractKeyFromPassphraseString("verysecretkey4"),
	}
	pubKeys := []string{
		privKeys[0].PublicKeyHex(),
		privKeys[1].PublicKeyHex(),
		privKeys[2].PublicKeyHex(),
		privKeys[3].PublicKeyHex(),
	}

	signers := strings.Join(pubKeys, ",")

	restClient := resty.New()

	rawNewDtrustAddressResponse, _ := restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
			"labels":					dtrustAddressLabel,
			"public_keys":				signers,
			"required_signatures":		"3",
			"address_type":				"witness_v0",
		}).Post("https://block.io/api/v2/get_new_dtrust_address?api_key=" + apiKey)

	parsedRes, _ := blockio.ParseResponse(rawNewDtrustAddressResponse.String())

	dtrustAddress = parsedRes.Data.(map[string]interface{})["address"].(string)
	fmt.Println("Our dTrust Address: " + dtrustAddress)

	rawWithdrawFromLabelsRes, _ := restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
			"from_labels":	"default",
			"to_address":	dtrustAddress,
			"amounts":		"0.1",
		}).Post("https://block.io/api/v2/withdraw_from_labels?api_key=" + apiKey)


	withdrawData, _ := blockio.ParseSignatureResponse(rawWithdrawFromLabelsRes.String())
	signatureReq, _ := blockio.SignWithdrawRequest(pin, withdrawData)

	signAndFinalizeWithdrawRes, _ := restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
			"signature_data": signatureReq,
		}).Post("https://block.io/api/v2/sign_and_finalize_withdrawal?api_key=" + apiKey)

	fmt.Println("Withdrawal Response: ")
	fmt.Println(blockio.ParseResponse(signAndFinalizeWithdrawRes.String()))

	rawDtrustAddressBalance, _ := restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
			"address":	dtrustAddress,
		}).Post("https://block.io/api/v2/get_dtrust_address_balance?api_key=" + apiKey)

	fmt.Println("Dtrust address balance: ")
	fmt.Println(blockio.ParseResponse(rawDtrustAddressBalance.String()))

	rawDefaultAddressRes, _ := restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
			"label":	"default",
		}).Post("https://block.io/api/v2/get_address_by_label?api_key=" + apiKey)

	parsedRes, _ = blockio.ParseResponse(rawDefaultAddressRes.String())
	normalAddress := parsedRes.Data.(map[string]interface{})["address"].(string)

	fmt.Println("Withdrawing from dtrust_address_label to the 'default' label in normal multisig")

	rawWithdrawFromDtrustRes, _ := restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
			"from_address":	dtrustAddress,
			"to_address":	normalAddress,
			"amounts":		"0.01",
		}).Post("https://block.io/api/v2/withdraw_from_dtrust_address?api_key=" + apiKey)

	unsignedSignatureReq, _ := blockio.ParseSignatureResponse(rawWithdrawFromDtrustRes.String())
	fmt.Println("Withdraw from Dtrust Address response:");
	fmt.Println(unsignedSignatureReq)

	signatureReq, _ = blockio.SignDtrustRequest(privKeys, unsignedSignatureReq)

	fmt.Println("Our Signed Request: ")
	fmt.Println(signatureReq)

	signAndFinalizeRes, _ := restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
			"signature_data": signatureReq,
		}).Post("https://block.io/api/v2/sign_and_finalize_withdrawal?api_key=" + apiKey)

	parsedSignAndFinalizeRes, _ := blockio.ParseResponse(signAndFinalizeRes.String())

	fmt.Println("Finalize Withdrawal: ");
	fmt.Println(parsedSignAndFinalizeRes.Data)

	rawDtrustTransactions, _ := restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
			"address":	dtrustAddress,
			"type":		"sent",
		}).Post("https://block.io/api/v2/get_dtrust_transactions?api_key=" + apiKey)

	fmt.Println("Get transactions sent by our dtrust_address_label address: ")
	fmt.Println(rawDtrustTransactions.String())
}