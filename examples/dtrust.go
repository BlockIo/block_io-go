package examples

import (
	"encoding/json"
	"fmt"
	"github.com/BlockIo/block_io-go/lib"
	"github.com/BlockIo/block_io-go/lib/BlockIo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

func dtrust() BlockIo.Client {
	var blockIo BlockIo.Client
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	api_key := os.Getenv("API_KEY")
	pin := os.Getenv("PIN")
	blockIo.Instantiate(api_key,pin,2, BlockIo.Options{})

	return blockIo
}

func RunDtrustExample() {
	blockIo := dtrust()
	dtrustAddress := ""
	dtrustAddressLabel := "dTrust1_witness_v0"

	privKeys := []string{
		lib.ExtractKeyFromPassphraseString("verysecretkey1"),
		lib.ExtractKeyFromPassphraseString("verysecretkey2"),
		lib.ExtractKeyFromPassphraseString("verysecretkey3"),
		lib.ExtractKeyFromPassphraseString("verysecretkey4"),
	}
	pubKeys := []string{
		lib.ExtractPubKeyFromPassphraseString("verysecretkey1"),
		lib.ExtractPubKeyFromPassphraseString("verysecretkey2"),
		lib.ExtractPubKeyFromPassphraseString("verysecretkey3"),
		lib.ExtractPubKeyFromPassphraseString("verysecretkey4"),
	}

	signers := strings.Join(pubKeys, ",")
	res := blockIo.GetNewDtrustAddress("{\"label\": \"" + dtrustAddressLabel + "\", \"public_keys\": \"" + signers + "\", \"required_signatures\": \"3\", \"address_type\": \"witness_v0\"}")
	if res["error_message"] != nil {
		fmt.Println("Error: ", res["error_message"])

		res = nil
		res = blockIo.GetDtrustAddressByLabel("{\"label\": \"" + dtrustAddressLabel + "\"}")
	}
	dtrustAddress = fmt.Sprintf("%v", res["address"])
	fmt.Println("Our dTrust Address:", dtrustAddress)
	res = nil
	res = blockIo.WithdrawFromLabels("{\"from_labels\": \"default\", \"to_address\": \"" + dtrustAddress + "\", \"amounts\": \"0.001\"}");
	fmt.Println("Withdrawal Response:", res)
	res = nil
	res = blockIo.GetDtrustAddressBalance("{\"label\": \"" + dtrustAddressLabel + "\"}")
	fmt.Println("Dtrust address label Balance: ", res);
	res = nil
	res = blockIo.GetAddressByLabel("{\"label\": \"default\"}")
	normalAddress := fmt.Sprintf("%v", res["address"])
	fmt.Println("Withdrawing from dtrust_address_label to the 'default' label in normal multisig")
	res = blockIo.WithdrawFromDtrustAddress("{\"from_labels\": \"" + dtrustAddressLabel + "\", \"to_address\": \"" + normalAddress + "\", \"amounts\": \"0.0009\"}");
	fmt.Println("Withdraw from Dtrust Address response: ", res)
	jsonString, _ := json.Marshal(res)
	var pojo BlockIo.SignatureData
	err := json.Unmarshal([]byte(string(jsonString)), &pojo)
	if err != nil {
		fmt.Println("Invalid response conversion:", err)
		return
	}

	for i := 0; i < len(pojo.Inputs); i++ {
		for j := 0; j < len(pojo.Inputs[i].Signers); j++ {
			pojo.Inputs[i].Signers[j].SignedData = lib.SignInputs(privKeys[j],pojo.Inputs[i].DataToSign)
		}
	}
	pojoMarshalled, _ := json.Marshal(pojo)

	fmt.Println("Our Signed Request:", string(pojoMarshalled))
	fmt.Println("Finalize Withdrawal:" )

	fmt.Println(blockIo.SignAndFinalizeWithdrawal(string(pojoMarshalled)))
	fmt.Println("Get transactions sent by our dtrust_address_label address: ")
	fmt.Println(blockIo.GetDtrustTransactions("{\"type\": \"sent\", \"labels\": \"" + dtrustAddressLabel + "\"}"))
}