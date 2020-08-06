package examples

import (
	"fmt"
	"github.com/BlockIo/block_io-go/lib/BlockIo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func maxWithdrawal() BlockIo.Client {
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

func RunMaxWithdrawalExample() {
	blockIo := maxWithdrawal()

	var balance string = blockIo.GetBalance(nil)["available_balance"].(string)

	fmt.Println("Balance:", balance)

	for {
		res := blockIo.Withdraw(map[string]interface{}{
			"to_address":os.Getenv("TO_ADDRESS"),
			"amount": balance,
		})
		if res["reference_id"] != nil { fmt.Println(res) }


		maxWithdrawString := res["max_withdrawal_available"].(string)

		fmt.Println("Max Withdraw Available:", maxWithdrawString)

		maxWithdraw,_ := strconv.ParseFloat(maxWithdrawString,64)

		if maxWithdraw == 0 { break }

		blockIo.Withdraw(map[string]interface{}{
			"to_address":os.Getenv("TO_ADDRESS"),
			"amount":maxWithdrawString,
		})

		break
	}

	balance = blockIo.GetBalance(nil)["available_balance"].(string)
	fmt.Println("Final Balance:", balance)
}
