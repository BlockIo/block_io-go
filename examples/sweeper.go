package examples

import (
	"fmt"
	"github.com/BlockIo/block_io-go/lib/BlockIo"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func sweeper() BlockIo.Client {
	var blockIo BlockIo.Client
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	api_key := os.Getenv("API_KEY")
	blockIo.Instantiate(api_key,"",2,"")

	return blockIo
}

func RunSweeperExample() {
	blockIo := sweeper()

	if (os.Getenv("TO_ADDRESS") == "") ||
			(os.Getenv("PRIVATE_KEY_FROM_ADDRESS") == "") ||
			(os.Getenv("FROM_ADDRESS") == "") {

		fmt.Println("Error: Missing parameters from env.")
		return
	}

	res := blockIo.SweepFromAddress("{" +
		"\"to_address\": \""          +     os.Getenv("TO_ADDRESS") +
		"\", \"private_key\": \""     +     os.Getenv("PRIVATE_KEY_FROM_ADDRESS") +
		"\", \"from_addresss\": \""   +     os.Getenv("FROM_ADDRESS") +
		"\"}")

	fmt.Println("Sweep Res:",res)
}
