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
	fmt.Println("Get New Address: ", blockIo.GetNewAddress(map[string]interface{}{
		"label":"testDest15",
	}))
	fmt.Println("Withdraw from labels: ", blockIo.WithdrawFromLabels(map[string]interface{}{
		"from_labels": "default",
		"to_label": "testDest15",
		"amount":"2.5",
	}))
	fmt.Println("Get Address Balance: ", blockIo.GetAddressBalance(map[string]interface{}{
		"labels": "default, testDest15",
	}))
	fmt.Println("Get Sent Transactions: ", blockIo.GetTransactions(map[string]interface{}{
		"type":"sent",
	}))
	fmt.Println("Get Received Transactions: ", blockIo.GetTransactions(map[string]interface{}{
		"type":"received",
	}))
	fmt.Println("Get Current Price: ", blockIo.GetCurrentPrice(map[string]interface{}{
		"base_price":"BTC",
	}))
}
