package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/go-resty/resty/v2"
	"os"
)

func main(){
	godotenv.Load(".env")
	apiKey := os.Getenv("API_KEY")
	restClient := resty.New()

	rawGetAddressRes, _ := restClient.R().
		Get("https://block.io/api/v2/get_new_address?api_key=" + apiKey)

	fmt.Println("get_new_address response:")
	fmt.Print(rawGetAddressRes)
}
