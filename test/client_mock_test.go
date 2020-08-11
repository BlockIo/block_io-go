package test

import (
	"bytes"
	"encoding/json"
	"github.com/BlockIo/block_io-go/lib/BlockIo"
	"github.com/jarcoal/httpmock"
	"io/ioutil"
	"net/http"
	"os"

	"fmt"
	"testing"
)
var blockIo BlockIo.Client
var api_key string

var signAndFinalizeWithdraw map[string]interface{}
var signAndFinalizeSweep map[string]interface{}
var signAndFinalizeDtrust map[string]interface{}

var withdrawRequestBodyContent map[string]interface{}
var dTrustRequestBodyContent map[string]interface{}
var sweepRequestBodyContent map[string]interface{}

func setup(pin string){
	blockIo = basic(pin)
	api_key = "0000-0000-0000-0000"
	httpmock.ActivateNonDefault(blockIo.RestClient.GetClient())

	withdrawRequestBodyContent = map[string]interface{}{
		"from_labels": "testdest",
		"amounts": "100",
		"to_labels": "default"}
	dTrustRequestBodyContent = map[string]interface{}{
		"to_addresses": "QhSWVppS12Fqv6dh3rAyoB18jXh5mB1hoC",
		"from_address": "tltc1q8y9naxlsw7xay4jesqshnpeuc0ap8fg9ejm2j2memwq4ng87dk3s88nr5j",
		"amounts": 0.09 }
	sweepRequestBodyContent = map[string]interface{}{
		"to_address": "QhSWVppS12Fqv6dh3rAyoB18jXh5mB1hoC",
		"from_address": "tltc1qpygwklc39wl9p0wvlm0p6x42sh9259xdjl059s",
		"private_key": "cTYLVcC17cYYoRjaBu15rEcD5WuDyowAw562q2F1ihcaomRJENu5"}

	readSignAndFinalizeWithdrawRequestJson()
	readSignAndFinalizeSweepRequestJson()
	readSignAndFinalizeDtrustRequestJson()
}

func readSignAndFinalizeWithdrawRequestJson(){
	jsonFile, err := os.Open("data/sign_and_finalize_withdrawal_request.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	byteValue = bytes.TrimPrefix(byteValue, []byte("\xef\xbb\xbf"))
	var result map[string]interface{}
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		fmt.Println(err)
	}
	signAndFinalizeWithdraw = result
}

func readSignAndFinalizeSweepRequestJson(){
	jsonFile, err := os.Open("data/sign_and_finalize_sweep_request.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	byteValue = bytes.TrimPrefix(byteValue, []byte("\xef\xbb\xbf"))
	var result map[string]interface{}
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		fmt.Println(err)
	}
	signAndFinalizeSweep = result
}

func readSignAndFinalizeDtrustRequestJson(){
	jsonFile, err := os.Open("data/sign_and_finalize_dtrust_withdrawal_request.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	byteValue = bytes.TrimPrefix(byteValue, []byte("\xef\xbb\xbf"))
	var result map[string]interface{}
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		fmt.Println(err)
	}
	signAndFinalizeDtrust = result
}

func setupDtrustStub() {
	httpmock.RegisterResponder("POST", "https://block.io/api/v2/withdraw_from_dtrust_address",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, fileToString("data/withdraw_from_dtrust_address_response.json"))
			if err != nil {
				fmt.Println("err",err)
			}
			return resp, err
		},
	)

	httpmock.RegisterResponder("POST", "https://block.io/api/v2/sign_and_finalize_withdrawal",
		func(req *http.Request) (*http.Response, error) {
			responseBody := make(map[string]interface{})
			if err := json.NewDecoder(req.Body).Decode(&responseBody); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			var pojo BlockIo.SignatureData
			err := json.Unmarshal([]byte(responseBody["signature_data"].(string)), &pojo)

			var pojo2 BlockIo.SignatureData
			err = json.Unmarshal([]byte(signAndFinalizeDtrust["signature_data"].(string)), &pojo2)
			check := false

			if pojo.ReferenceID != pojo2.ReferenceID {
				check = true
			}
			if len(pojo.Inputs) != len(pojo2.Inputs){

				check = true
			}
			for i := 0; i < len(pojo.Inputs); i++ {
				if pojo.Inputs[i].DataToSign != pojo2.Inputs[i].DataToSign || pojo.Inputs[i].InputNo != pojo2.Inputs[i].InputNo ||
					pojo.Inputs[i].SignaturesNeeded != pojo2.Inputs[i].SignaturesNeeded {
					check = true
				}
				for j := 0; j < len(pojo.Inputs[i].Signers); j++ {
					if pojo.Inputs[i].Signers[j]!=pojo2.Inputs[i].Signers[j] {
						check = true
					}
				}
			}
			var resp *http.Response
			if check == true {
				resp, err = httpmock.NewJsonResponse(500, "Error")
			} else {
				resp, err = httpmock.NewJsonResponse(200, fileToString("data/success_response.json"))
			}

			if err != nil {
				fmt.Println("err",err)
			}
			return resp, err
		},
	)

}

func setupWithdrawStub() {
	httpmock.RegisterResponder("POST", "https://block.io/api/v2/withdraw",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, fileToString("data/withdraw_response.json"))
			if err != nil {
				fmt.Println("err",err)
			}
			return resp, err
		},
	)

	httpmock.RegisterResponder("POST", "https://block.io/api/v2/sign_and_finalize_withdrawal",
		func(req *http.Request) (*http.Response, error) {

			responseBody := make(map[string]interface{})
			if err := json.NewDecoder(req.Body).Decode(&responseBody); err != nil {
				return httpmock.NewStringResponse(400, "error parsing request body"), nil
			}
			var pojo BlockIo.SignatureData
			err := json.Unmarshal([]byte(responseBody["signature_data"].(string)), &pojo)

			var pojo2 BlockIo.SignatureData
			err = json.Unmarshal([]byte(signAndFinalizeWithdraw["signature_data"].(string)), &pojo2)

			check := false

			if pojo.ReferenceID != pojo2.ReferenceID {
				check = true
			}
			if len(pojo.Inputs) != len(pojo2.Inputs){
				check = true
			}
			for i := 0; i < len(pojo.Inputs); i++ {
				if pojo.Inputs[i].DataToSign != pojo2.Inputs[i].DataToSign || pojo.Inputs[i].InputNo != pojo2.Inputs[i].InputNo ||
					pojo.Inputs[i].SignaturesNeeded != pojo2.Inputs[i].SignaturesNeeded {
					check = true
				}
				for j := 0; j < len(pojo.Inputs[i].Signers); j++ {
					if pojo.Inputs[i].Signers[j]!=pojo2.Inputs[i].Signers[j] {
						check = true
					}
				}
			}
			var resp *http.Response
			if check == true {
				resp, err = httpmock.NewJsonResponse(500, "Error")
			} else {
				resp, err = httpmock.NewJsonResponse(200, fileToString("data/success_response.json"))
			}

			if err != nil {
				fmt.Println("err",err)
			}
			return resp, err
		},
	)
}


func setupSweepStub() {
	httpmock.RegisterResponder("POST", "https://block.io/api/v2/sweep_from_address",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, fileToString("data/sweep_from_address_response.json"))
			if err != nil {
				fmt.Println("err",err)
			}
			return resp, err
		},
	)

	httpmock.RegisterResponder("POST", "https://block.io/api/v2/sign_and_finalize_sweep",
		func(req *http.Request) (*http.Response, error) {
			responseBody := make(map[string]interface{})
			if err := json.NewDecoder(req.Body).Decode(&responseBody); err != nil {
				return httpmock.NewStringResponse(400, "error parsing request body"), nil
			}
			var pojo BlockIo.SignatureData
			err := json.Unmarshal([]byte(responseBody["signature_data"].(string)), &pojo)

			var pojo2 BlockIo.SignatureData
			err = json.Unmarshal([]byte(signAndFinalizeSweep["signature_data"].(string)), &pojo2)
			check := false

			if pojo.ReferenceID != pojo2.ReferenceID {
				fmt.Println(1)
				check = true
			}
			if len(pojo.Inputs) != len(pojo2.Inputs){
				fmt.Println(2)
				check = true
			}
			for i := 0; i < len(pojo.Inputs); i++ {
				if pojo.Inputs[i].DataToSign != pojo2.Inputs[i].DataToSign || pojo.Inputs[i].InputNo != pojo2.Inputs[i].InputNo ||
					pojo.Inputs[i].SignaturesNeeded != pojo2.Inputs[i].SignaturesNeeded {
					fmt.Println(3)
					check = true
				}
				for j := 0; j < len(pojo.Inputs[i].Signers); j++ {
					if pojo.Inputs[i].Signers[j]!=pojo2.Inputs[i].Signers[j] {
						fmt.Println(pojo.Inputs[i].Signers[j].SignedData,pojo2.Inputs[i].Signers[j].SignedData)
						check = true
					}
				}
			}
			var resp *http.Response
			if check == true {
				resp, err = httpmock.NewJsonResponse(500, "Error")
			} else {
				resp, err = httpmock.NewJsonResponse(200, fileToString("data/success_response.json"))
			}

			if err != nil {
				fmt.Println("err",err)
			}
			return resp, err
		},
	)
}

func basic(pin string) BlockIo.Client {
	var blockIo BlockIo.Client

	api_key := "0000-0000-0000-0000"
	blockIo.Instantiate(api_key,pin,2, BlockIo.Options{})

	return blockIo
}

func fileToString(path string) map[string]interface{} {
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	byteValue = bytes.TrimPrefix(byteValue, []byte("\xef\xbb\xbf"))
	var result map[string]interface{}
	err = json.Unmarshal([]byte(byteValue), &result)
	return result
}



func TestWithdraw(t *testing.T){
	defer httpmock.DeactivateAndReset()
	setup("blockiotestpininsecure")
	setupWithdrawStub()

	resp, err := blockIo.Withdraw(withdrawRequestBodyContent)
	if err != nil {
		t.Error(err)
	}
	if resp["txid"]==nil {
		t.Error("Withdraw Response not correct")
	}
}

func TestSweep(t *testing.T){
	defer httpmock.DeactivateAndReset()
	setup("")
	setupSweepStub()

	resp, err := blockIo.SweepFromAddress(sweepRequestBodyContent)
	if err != nil {
		t.Error(err)
	}
	if resp["txid"]==nil {
		t.Error("Withdraw Response not correct")
	}
}

func TestDtrust(t *testing.T){
	defer httpmock.DeactivateAndReset()
	setup("blockiotestpininsecure")
	setupDtrustStub()

	resp, err := blockIo.WithdrawFromDtrustAddress(sweepRequestBodyContent)
	if err != nil {
		t.Error(err)
	}
	if resp==nil {
		t.Error("Withdraw Response not correct")
	}
}

