package BlockIo

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"strconv"
	"strings"

	"github.com/BlockIo/block_io-go/lib"
)
type Options struct {
	AllowNoPin bool
	ApiUrl     string
}
type Client struct {
	options        Options
	apiUrl         string
	pin            string
	aesKey         string
	apiKey         string
	version        int
	server         string
	port           string
	restClient     *resty.Client
}

const defaultVersion = 2
const defaultServer = ""
const defaultPort = ""
const host = "block.io"

func (blockIo *Client) Instantiate(apiKey string, pin string, version int, opts Options) {

	blockIo.options = opts
	blockIo.apiUrl = blockIo.options.ApiUrl

	blockIo.pin = pin
	blockIo.aesKey = ""

	blockIo.apiKey = apiKey
	if version == -1 {
		blockIo.version = defaultVersion
	} else {
		blockIo.version = version
	}
	blockIo.server = defaultServer
	blockIo.port = defaultPort
	if blockIo.pin != "" {
		blockIo.aesKey = lib.PinToAesKey(blockIo.pin)
	}

	serverString := blockIo.server
	if serverString != "" {
		serverString = serverString + "."
	}

	portString := ""
	if blockIo.port != "" { portString = ":" + blockIo.port }

	if blockIo.apiUrl == "" {
		blockIo.apiUrl = "https://" + serverString + host + portString + "/api/v" + strconv.Itoa(blockIo.version) + "/"
	}
	blockIo.restClient = resty.New()
}

func (blockIo *Client) get(path string) (string, error) {
	client := resty.New()

	resp, err := client.R().
		EnableTrace().
		Get(blockIo.constructUrl(path))

	return resp.String(), err

}
func (blockIo *Client) post(jsonInput string, path string) (string, error) {
	var argsObj map[string]interface{}
	parseErr := json.Unmarshal([]byte(jsonInput), &argsObj)

	if parseErr != nil {
		return "", parseErr
	}

	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetBody(argsObj).
		Post(blockIo.constructUrl(path))

	return resp.String(), err
}

func (blockIo *Client) _withdraw(Method string, Path string, args map[string]interface{}) (map[string]interface{}, error) {

	var pin string
	if args["pin"] != nil {
		pin = args["pin"].(string)
	} else {
		pin = blockIo.pin
	}
	args["pin"] = ""
	argsObj, argsErr := json.Marshal(args)
	if argsErr != nil {
		return map[string]interface{}{}, argsErr
	}
	res, err := blockIo._request(Method, Path, string(argsObj))
	if err != nil {
		return map[string]interface{}{}, err
	}
	jsonString, resErr := json.Marshal(res)
	if resErr != nil {
		return map[string]interface{}{}, resErr
	}
	var pojo SignatureData
	err = json.Unmarshal(jsonString, &pojo)
	if err != nil {
		return map[string]interface{}{}, err
	}
	if (pojo.ReferenceID == "" || pojo.EncryptedPassphrase == EncryptedPassphrase{} ||
		pojo.EncryptedPassphrase.Passphrase == "") {
		return res, nil
	}
	if pin == "" {
		if blockIo.options.AllowNoPin {
			return res, nil
		}
	}
	var encryptedPassphrase = pojo.EncryptedPassphrase.Passphrase
	var aesKey string
	if blockIo.aesKey != "" { aesKey = blockIo.aesKey
	} else {aesKey = lib.PinToAesKey(pin) }
	privKey := lib.ExtractKeyFromEncryptedPassphrase(encryptedPassphrase,aesKey)
	pubKey := lib.ExtractPubKeyFromEncryptedPassphrase(encryptedPassphrase,aesKey)
	if pubKey != pojo.EncryptedPassphrase.SignerPublicKey {
		return map[string]interface{}{}, errors.New("Public key mismatch. Invalid Secret PIN detected.")
	}

	for i := 0; i < len(pojo.Inputs); i++ {
		for j := 0; j < len(pojo.Inputs[i].Signers); j++ {
			pojo.Inputs[i].Signers[j].SignedData = lib.SignInputs(privKey,pojo.Inputs[i].DataToSign)
		}
	}
	pojo.EncryptedPassphrase = EncryptedPassphrase{}
	pojoMarshalled, pojoErr := json.Marshal(pojo)
	if pojoErr != nil {
		return map[string]interface{}{}, pojoErr
	}
	return blockIo._request(Method,"sign_and_finalize_withdrawal",string(pojoMarshalled))
}

func (blockIo *Client) _sweep(Method string, Path string, args  map[string]interface{}) (map[string]interface{}, error) {

	if args["to_address"] == nil {
		return map[string]interface{}{}, errors.New("Missing mandatory private_key argument.")
	}

	privKeyStr := args["private_key"].(string)
	keyFromWif, _ := lib.FromWIF(privKeyStr)
	args["public_key"] = lib.PubKeyFromWIF(privKeyStr)
	args["private_key"] = ""

	argsObjMarshalled, argsErr := json.Marshal(args)
	if argsErr != nil {
		return map[string]interface{}{}, argsErr
	}

	res, err := blockIo._request(Method, Path, string(argsObjMarshalled))
	if err != nil {
		return map[string]interface{}{}, err
	}
	jsonString, resErr := json.Marshal(res)

	if resErr != nil {
		return map[string]interface{}{}, resErr
	}
	var pojo SignatureData
	err = json.Unmarshal([]byte(string(jsonString)), &pojo)
	if err != nil {
		return map[string]interface{}{}, err
	}
	if pojo.ReferenceID == "" {
		return map[string]interface{}{}, errors.New("Empty Reference ID")
	}

	for i := 0; i < len(pojo.Inputs); i++ {
		for j := 0; j < len(pojo.Inputs[i].Signers); j++ {
			pojo.Inputs[i].Signers[j].SignedData = lib.SignInputs(keyFromWif,pojo.Inputs[i].DataToSign)
		}
	}
	pojo.EncryptedPassphrase = EncryptedPassphrase{}
	pojoMarshalled, pojoErr := json.Marshal(pojo)
	if pojoErr != nil {
		return map[string]interface{}{}, pojoErr
	}
	return blockIo._request(Method,"sign_and_finalize_sweep",string(pojoMarshalled))
}

func (blockIo *Client) _request(Method string, Path string, args string) (map[string]interface{}, error) {
	var res map[string]interface{}
	if Method == "POST" {
		if strings.Contains(Path, "sign_and_finalize") {

			var postObj = map[string]string{"signature_data": args}

			temp, postMarshalErr := json.Marshal(postObj)
			if postMarshalErr != nil {
				return map[string]interface{}{}, postMarshalErr
			}
			args = string(temp)
		}
		resString, postErr := blockIo.post(args, Path)
		if postErr != nil {
			return map[string]interface{}{}, postErr
		}
		marshalErr := json.Unmarshal([]byte(resString), &res)
		if marshalErr != nil {
			return map[string]interface{}{}, marshalErr
		}
	} else {
		resString, getErr := blockIo.get(Path)
		if getErr != nil {
			return map[string]interface{}{}, getErr
		}
		marshalErr := json.Unmarshal([]byte(resString), &res)
		if marshalErr != nil {
			return map[string]interface{}{}, marshalErr
		}
	}
	jsonString, dataErr := json.Marshal(res["data"])
	if dataErr != nil {
		return map[string]interface{}{}, dataErr
	}
	res = nil
	err := json.Unmarshal(jsonString, &res)
	if err != nil {
		return map[string]interface{}{}, err
	}
	if res == nil {
		return map[string]interface{}{}, errors.New("No response from API server")
	}
	return res, nil
}

func (blockIo *Client) constructUrl(path string) string {
	return blockIo.apiUrl + path + "?api_key=" + blockIo.apiKey
}

