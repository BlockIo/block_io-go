package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"os"
)

func main(){
	godotenv.Load(".env")
	apiKey := os.Getenv("API_KEY")
	address := os.Getenv("ADDRESS")
	restClient := resty.New()

	rawGetAddressBalanceRes, _ := restClient.R().
		Get("https://block.io/api/v2/get_address_balance?api_key=" + apiKey + "&address=" + address)

	fmt.Println("get_address_balance response:")
	fmt.Print(rawGetAddressBalanceRes)
}
