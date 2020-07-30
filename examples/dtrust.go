package examples

import (
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
	blockIo.Instantiate(api_key,pin,2,"")

	return blockIo
}

func RunDtrustExample() {
	blockIo := dtrust()
	//dtrustAddress := nil
	dtrustAddressLabel := "dTrust1_witness_v0"

	//privKeys := []string{
	//	lib.ExtractKeyFromPassphraseString("verysecretkey1"),
	//	lib.ExtractKeyFromPassphraseString("verysecretkey2"),
	//	lib.ExtractKeyFromPassphraseString("verysecretkey3"),
	//	lib.ExtractKeyFromPassphraseString("verysecretkey4"),
	//}
	pubKeys := []string{
		lib.ExtractPubKeyFromPassphraseString("verysecretkey1"),
		lib.ExtractPubKeyFromPassphraseString("verysecretkey2"),
		lib.ExtractPubKeyFromPassphraseString("verysecretkey3"),
		lib.ExtractPubKeyFromPassphraseString("verysecretkey4"),
	}

	signers := strings.Join(pubKeys, ",")
	res := blockIo.GetNewDtrustAddress("{\"label\": \"" + dtrustAddressLabel + "\", \"public_keys\": \"" + signers + "\", \"required_signatures\": \"3\", \"address_type\": \"witness_v0\"}")
	fmt.Println(res)
	if res["reference_id"] == "" {

	}
}