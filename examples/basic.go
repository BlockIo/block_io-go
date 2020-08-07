package examples

import (
	"fmt"
	"github.com/BlockIo/block_io-go/lib/BlockIo"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func basic() BlockIo.Client {
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

func RunBasicExample(){
	blockIo := basic()
	res, _ := blockIo.GetNewAddress(map[string]interface{}{
		"label":"testDest15",
	})
	fmt.Println("Get New Address: ", res)
	res, _ = blockIo.WithdrawFromLabels(map[string]interface{}{
		"from_labels": "default",
		"to_label": "testDest15",
		"amount":"2.5",
	})
	fmt.Println("Withdraw from labels: ", res)
	res, _ = blockIo.GetAddressBalance(map[string]interface{}{
		"labels": "default, testDest15",
	})
	fmt.Println("Get Address Balance: ", res)
	res, _ = blockIo.GetTransactions(map[string]interface{}{
		"type":"sent",
	})
	fmt.Println("Get Sent Transactions: ", res)
	res, _ = blockIo.GetTransactions(map[string]interface{}{
		"type":"received",
	})
	fmt.Println("Get Received Transactions: ", res)
	res, _ = blockIo.GetCurrentPrice(map[string]interface{}{
		"base_price":"BTC",
	})
	fmt.Println("Get Current Price: ", res)
}
