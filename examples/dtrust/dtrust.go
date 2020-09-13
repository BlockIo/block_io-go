package main

import (
	"fmt"
	blockio "github.com/BlockIo/block_io-go"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"os"
	"strings"
)

func main(){

	//load vars from .env
	godotenv.Load(".env")

	//load environment vars
	apiKey := os.Getenv("API_KEY")
	pin := os.Getenv("PIN")

	if (apiKey == "" || pin == "") {
		log.Fatal("Dtrust requires environment variables API_KEY and PIN")
	}

	dtrustAddress := ""

	//create a random address label
	dtrustAddressLabel := "dTrust1" + fmt.Sprintf("%f", rand.ExpFloat64())

	// create the private key objects for each private key
	// NOTE: in production environments you'll do this elsewhere
	privKeys := []*blockio.ECKey{
		blockio.ExtractKeyFromPassphraseString("verysecretkey1"),
		blockio.ExtractKeyFromPassphraseString("verysecretkey2"),
		blockio.ExtractKeyFromPassphraseString("verysecretkey3"),
	}

	// populate our pubkeys array from the keys we just generated
	// pubkey entries are expected in hexadecimal format
	pubKeys := []string{
		privKeys[0].PublicKeyHex(),
		privKeys[1].PublicKeyHex(),
		privKeys[2].PublicKeyHex(),
	}

	signers := strings.Join(pubKeys, ",")

	// instantiate a rest client
	restClient := resty.New()

	// create a dTrust address that requires 4 out of 5 keys (4 of ours, 1 at Block.io).
	// Block.io automatically adds +1 to specified required signatures because of its own key
	newDtrustAddrRes, newDtrustAddrErr := restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
			"labels":					dtrustAddressLabel,
			"public_keys":				signers,
			"required_signatures":		"3", // required signatures out of the set of signatures that we specified
			"address_type":				"witness_v0",
		}).Post("https://block.io/api/v2/get_new_dtrust_address?api_key=" + apiKey)

	if newDtrustAddrErr != nil {
		log.Fatal(newDtrustAddrErr)
	}

	parsedAddr, parseErr := ParseAddrResponse(newDtrustAddrRes.String())

	if parseErr != nil {
		log.Fatal(parseErr)
	}
	dtrustAddress = parsedAddr
	fmt.Println("Our dTrust Address: " + dtrustAddress)

	// let's send some coins to our new address
	withdrawFromLabelRes, withdrawFromLabelErr := restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
			"from_labels":	"default",
			"to_address":	dtrustAddress,
			"amounts":		"0.1",
		}).Post("https://block.io/api/v2/withdraw_from_labels?api_key=" + apiKey)

	if withdrawFromLabelErr != nil {
		log.Fatal(withdrawFromLabelErr)
	}

	signatureReq, signatureReqErr := blockio.SignWithdrawRequestJson(pin, withdrawFromLabelRes.String())

	if signatureReqErr != nil {
		log.Fatal(signatureReqErr)
	}

	signAndFinalizeWithdrawRes, signAndFinalizeWithdrawErr := restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
			"signature_data": signatureReq,
		}).Post("https://block.io/api/v2/sign_and_finalize_withdrawal?api_key=" + apiKey)

	if signAndFinalizeWithdrawErr != nil {
		log.Fatal(signAndFinalizeWithdrawErr)
	}

	fmt.Println("Withdrawal Response: ")
	fmt.Println(blockio.ParseResponse(signAndFinalizeWithdrawRes.String()))

	// check if some balance got there
	dtrustAddrBalanceRes, dtrustAddrBalErr := restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
			"address":	dtrustAddress,
		}).Post("https://block.io/api/v2/get_dtrust_address_balance?api_key=" + apiKey)

	if dtrustAddrBalErr != nil {
		log.Fatal(dtrustAddrBalErr)
	}

	fmt.Println("Dtrust address balance: ")
	fmt.Println(blockio.ParseResponse(dtrustAddrBalanceRes.String()))

	// find our non-dtrust default address so we can send coins back to it
	defaultAddrRes, defaultAddrErr := restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
			"label":	"default",
		}).Post("https://block.io/api/v2/get_address_by_label?api_key=" + apiKey)

	if defaultAddrErr != nil {
		log.Fatal(defaultAddrErr)
	}

	parsedAddr, parseErr = ParseAddrResponse(defaultAddrRes.String())

	if parseErr != nil {
		log.Fatal(parseErr)
	}

	normalAddress := parsedAddr

	fmt.Println("Withdrawing from dtrust_address_label to the 'default' label in normal multisig")

	// let's send the coins back to the default address
	withdrawDtrustRes, withdrawDtrustErr := restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
			"from_address":	dtrustAddress,
			"to_address":	normalAddress,
			"amounts":		"0.01",
		}).Post("https://block.io/api/v2/withdraw_from_dtrust_address?api_key=" + apiKey)

	if withdrawDtrustErr != nil {
		log.Fatal(withdrawDtrustErr)
	}

	fmt.Println("Withdraw from Dtrust Address response:");
	fmt.Println(withdrawDtrustRes)

	// Sign request with one key
	signatureReq, signErr := blockio.SignDtrustRequestWithKey(privKeys[0], withdrawDtrustRes.String())

	if signErr != nil {
		log.Fatal(signErr)
	}

	fmt.Println("Our Request signed with a single key: ")
	fmt.Println(signatureReq)

	signAndFinalizeRes, signAndFinalizeErr := restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
			"signature_data": signatureReq,
		}).Post("https://block.io/api/v2/sign_and_finalize_withdrawal?api_key=" + apiKey)

	if signAndFinalizeErr != nil {
		log.Fatal(signAndFinalizeErr)
	}

	// Sign request with 2 keys
	signatureReq, signErr = blockio.SignDtrustRequestWithKeys(privKeys[1:3], withdrawDtrustRes.String())

	if signErr != nil {
		log.Fatal(signErr)
	}

	fmt.Println("Our Request signed with 2 keys: ")
	fmt.Println(signatureReq)

	signAndFinalizeRes, signAndFinalizeErr = restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
			"signature_data": signatureReq,
		}).Post("https://block.io/api/v2/sign_and_finalize_withdrawal?api_key=" + apiKey)

	if signAndFinalizeErr != nil {
		log.Fatal(signAndFinalizeErr)
	}

	parsedSignAndFinalizeRes, parsedFinalizedErr := blockio.ParseResponse(signAndFinalizeRes.String())

	if parsedFinalizedErr != nil {
		log.Fatal(parsedFinalizedErr)
	}

	fmt.Println("Finalize Withdrawal: ");
	fmt.Println(parsedSignAndFinalizeRes.Data)

	dtrustTransactions, dtrustTransactionsErr := restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(map[string]interface{}{
			"address":	dtrustAddress,
			"type":		"sent",
		}).Post("https://block.io/api/v2/get_dtrust_transactions?api_key=" + apiKey)

	if dtrustTransactionsErr != nil {
		log.Fatal(dtrustTransactionsErr)
	}

	fmt.Println("Get transactions sent by our dtrust_address_label address: ")
	fmt.Println(dtrustTransactions.String())
}
