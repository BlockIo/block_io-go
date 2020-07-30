package main

import (
	"fmt"
	"github.com/BlockIo/block_io-go/lib/BlockIo"
)

func main() {

	var blockIo BlockIo.Client
	api_key := ""
	pin := ""
	blockIo.Instantiate(api_key,pin,2,"")
	res := blockIo.Withdraw("{\"from_labels\": \"shibe1\", \"to_label\": \"default\", \"amount\": \"0.01\"}")

	fmt.Println(res)
}
