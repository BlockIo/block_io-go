package main

import (
	"fmt"
	"github.com/BlockIo/block_io-go/lib/BlockIo"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var blockIo BlockIo.Client

	api_key := os.Getenv("API_KEY")
	pin := os.Getenv("PIN")
	blockIo.Instantiate(api_key,pin,2,"")
	res := blockIo.Withdraw("{\"to_addresses\": \"my9gXk65EzZUL962MSJadPXJFmJzPDc1WT\", \"amounts\": \"5.0\"}")

	res2 := blockIo.SweepFromAddress("{" +
		"\"to_address\": \""          +     "2N3ZkiF4WNfGAQYFRhivT3CLbUAYawbDNmV" +
		"\", \"private_key\": \""     +     "cUhedoiwPkprm99qfUKzixsrpN3w6wT2XrrMjqo3Yh1tHz8ykVKc" +
		"\", \"from_addresss\": \""   +     "my9gXk65EzZUL962MSJadPXJFmJzPDc1WT" +
		"\"}")

	fmt.Println(res,res2)
}
