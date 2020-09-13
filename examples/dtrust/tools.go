package main

import (
	"encoding/json"
	blockio "github.com/BlockIo/block_io-go"
)

type AddrResponse struct {
	Status string        `json:"status"`
	Data   Address	 	 `json:"data"`
}

type Address struct {
	Val string `json:"address"`
}

func ParseAddrResponse(res string) (string, error){
	resErr := blockio.ValidateResponseJson(res)
	if resErr != nil {
		return "", resErr
	}

	var addrRes AddrResponse
	marshalErr := json.Unmarshal([]byte(res), &addrRes)
	if marshalErr != nil {
		return "", marshalErr
	}

	return addrRes.Data.Val, nil
}