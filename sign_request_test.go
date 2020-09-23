package block_io_go

import (
	"io"
	"os"
	"strings"
	"testing"
)

var signingPin string
var withdrawReqJson string
var withdrawSigned string

func SignWithdrawRequestSetup(t *testing.T) {
	var reqBuf strings.Builder
	var signedBuf strings.Builder
	signingPin = "Was1qWas1q"

	withdrawReqFile, err := os.Open("data/withdraw_request.json")
	if err != nil {
		t.Error(err)
	}
	defer withdrawReqFile.Close()

	_, err = io.Copy(&reqBuf, withdrawReqFile)
	if err != nil {
		t.Error(err)
	}

	withdrawReqJson = reqBuf.String()

	withdrawSignedFile, err := os.Open("data/withdraw_signed.json")
	if err != nil {
		t.Error(err)
	}
	defer withdrawSignedFile.Close()

	_, err = io.Copy(&signedBuf, withdrawSignedFile)
	if err != nil {
		t.Error(err)
	}
	withdrawSigned = signedBuf.String()
}

func TestWithdraw(t *testing.T) {
	SignWithdrawRequestSetup(t)
	signatureReq, signErr := SignWithdrawRequestJson(signingPin, withdrawReqJson)
	if signErr != nil {
		t.Error(signErr)
	}
	if signatureReq != withdrawSigned {
		t.Error("SignWithdrawRequestJson failed")
	}
}
