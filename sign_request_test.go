package block_io_go

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
)

var signingPin string

var withdrawReqJson string
var sweepReqJson string
var dTrustReqJson string

var key *ECKey
var keys []*ECKey

func SignWithdrawRequestSetup(t *testing.T) {
	var reqBuf strings.Builder
	signingPin = "Was1qWas1q"

	withdrawReqFile, err := os.Open("fixtures/withdraw_request.json")
	if err != nil {
		t.Error(err)
	}
	defer withdrawReqFile.Close()

	_, err = io.Copy(&reqBuf, withdrawReqFile)
	if err != nil {
		t.Error(err)
	}

	withdrawReqJson = reqBuf.String()
}

func SignSweepRequestSetup(t *testing.T) {
	var reqBuf strings.Builder
	signingPin = "Was1qWas1q"

	sweepReqFile, err := os.Open("fixtures/sweep_request.json")
	if err != nil {
		t.Error(err)
	}
	defer sweepReqFile.Close()

	_, err = io.Copy(&reqBuf, sweepReqFile)
	if err != nil {
		t.Error(err)
	}

	sweepReqJson = reqBuf.String()

	ecKey, wifErr := FromWIF("cUhedoiwPkprm99qfUKzixsrpN3w6wT2XrrMjqo3Yh1tHz8ykVKc")
	if wifErr != nil {
		t.Errorf("Error extracting key from WIF: %s", wifErr)
	}
	key = ecKey
}

func SignDtrustRequestSetup(t *testing.T) {
	var reqBuf strings.Builder
	signingPin = "Was1qWas1q"

	dTrustReqFile, err := os.Open("fixtures/dtrust_request.json")
	if err != nil {
		t.Error(err)
	}
	defer dTrustReqFile.Close()

	_, err = io.Copy(&reqBuf, dTrustReqFile)
	if err != nil {
		t.Error(err)
	}

	dTrustReqJson = reqBuf.String()

	keys = []*ECKey{
		ExtractKeyFromPassphraseString("verysecretkey2"),
		ExtractKeyFromPassphraseString("verysecretkey3"),
	}
}

func ParseResult(str string) (*SignatureData, error) {
	var data SignatureData
	err := json.Unmarshal([]byte(str), &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func compareSignedRequest(sigs expSignedInputs, request *SignatureData, sigObj *SignatureData, t *testing.T) {

	// test if we got the correct amount of inputs
	if len(sigObj.Inputs) != len(sigs) {
		t.Errorf("Expected %d inputs, got %d", len(sigs), len(sigObj.Inputs))
	}

	// loop through all inputs
	for i := 0; i < len(sigObj.Inputs); i++ {

		// cache some vars
		indexedInput := sigObj.Inputs[i].InputNo
		expectedNumSigs := len(sigs[indexedInput])
		actualNumSigs := len(sigObj.Inputs[i].Signers)

		// check if the number of signatures matched
		if actualNumSigs != expectedNumSigs {
			t.Errorf("Expected %d signed inputs, got %d", expectedNumSigs, actualNumSigs)
		}

		// loop through each signer per input
		for j := 0; j < expectedNumSigs; j++ {

			origSignerObj := request.Inputs[i].Signers[j]
			actSignerObj := sigObj.Inputs[i].Signers[j]
			expSig := sigs[indexedInput][j]
			actSig := actSignerObj.SignedData

			if actSig == nil {
				actSig = "null"
			}
			// make sure the signatures are correct
			if expSig != actSig {
				t.Errorf(
					"Signature mismatch on input_no %d, signature %d\n  expected: %s\n  got:      %s",
					indexedInput, j, expSig, actSig)
			}

			// make sure the pubkey has not been changed
			if actSignerObj.SignerPublicKey != origSignerObj.SignerPublicKey {
				t.Errorf(
					"Public key mismatch on input_no %d, signature %d\n  expected: %s\n  got:      %s",
					indexedInput, j, origSignerObj.SignerPublicKey, actSignerObj.SignerPublicKey)
			}
		}
	}
}

type expSignedInputs map[int64][]string

func TestWithdraw(t *testing.T) {
	SignWithdrawRequestSetup(t)
	signatureReq, signErr := SignWithdrawRequestJson(signingPin, withdrawReqJson)
	if signErr != nil {
		t.Errorf("Signing threw an error: %s", signErr)
	}
	// signatures we expect
	expectedSigs := expSignedInputs{
		0: []string{"304502210084c918bb4d1bda7c8be9946bb5e4d073a992098effdc46e870a2c0bcb538774702204bac1b603ffaff3f744b3aa0494e743d17744e8490d990f3b458f5da6b08d29c"},
		1: []string{"3045022100be9d194a967a91c8f77db4c6bf0bd3d2fdb2235cd6a78328954c448e255aa17d02207bb89868300c594838b5fedad6f01810dbf9d2c97e44e267c63b121cab2dcdeb"},
	}
	reqObj, reqErr := ParseSignatureResponse(withdrawReqJson)
	if reqErr != nil {
		t.Errorf("Parsing the unsigned JSON threw an error: %s", reqErr)
	}
	// parse the JSON output - must output valid json and not throw error
	sigObj, parseErr := ParseResult(signatureReq)
	if parseErr != nil {
		t.Errorf("Parsing the signed JSON threw an error: %s", parseErr)
	}
	compareSignedRequest(expectedSigs, reqObj, sigObj, t)
}

func TestSweep(t *testing.T){
	SignSweepRequestSetup(t)
	signatureReq, signErr := SignRequestJsonWithKey(key, sweepReqJson)
	if signErr != nil {
		t.Errorf("Signing threw an error: %s", signErr)
	}
	// signatures we expect
	expectedSigs := expSignedInputs{
		0: []string{"30440220480a498bb662ec3993cf9e429143acb40d1b9c29e2d1256bbbb4c26506e9708902205f41f66c2ec5340f7d1dcbee1940b67faf768160c9c7e1ef3d199794bbc69d13"},
	}
	reqObj, reqErr := ParseSignatureResponse(sweepReqJson)
	if reqErr != nil {
		t.Errorf("Parsing the unsigned JSON threw an error: %s", reqErr)
	}
	// parse the JSON output - must output valid json and not throw error
	sigObj, parseErr := ParseResult(signatureReq)
	if parseErr != nil {
		t.Errorf("Parsing the signed JSON threw an error: %s", parseErr)
	}
	compareSignedRequest(expectedSigs, reqObj, sigObj, t)
}

func TestDtrust(t *testing.T){
	SignDtrustRequestSetup(t)
	signatureReq, signErr := SignRequestJsonWithKeys(keys, dTrustReqJson)
	if signErr != nil {
		t.Errorf("Signing threw an error: %s", signErr)
	}
	// signatures we expect
	expectedSigs := expSignedInputs{
		0: []string{
			"null",
			"3044022061fba610968257004b3668a294df0e4356f752e9916998988eabe5333ed5b7d702204e720d189b7cd4583dbe5a96e23f10117c9383401a1a1743dcf5a08fd8dbbfc5",
			"304502210085bb574741c747250d7ca09558583fd0831a774d8c3f57dd979a2dea47ad8f26022038faf15eeb88537647cb22c5300fefcde7cf7d35c3e041114af2feeddac3af25"},
	}
	reqObj, reqErr := ParseSignatureResponse(dTrustReqJson)
	if reqErr != nil {
		t.Errorf("Parsing the unsigned JSON threw an error: %s", reqErr)
	}
	// parse the JSON output - must output valid json and not throw error
	sigObj, parseErr := ParseResult(signatureReq)
	if parseErr != nil {
		t.Errorf("Parsing the signed JSON threw an error: %s", parseErr)
	}
	compareSignedRequest(expectedSigs, reqObj, sigObj, t)
}
