package BlockIo

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"strconv"
	"strings"

	"github.com/BlockIo/block_io-go/lib"
)
type Options struct {
	allowNoPin 	string
	apiUrl  	string
}
type Client struct {
	options        Options
	apiUrl         string
	pin            string
	aesKey         string
	apiKey         string
	version        int
	defaultVersion int
	server         string
	defaultServer  string
	port           string
	defaultPort    string
	host           string
	restClient     *resty.Client
}

func (blockIo *Client) Instantiate(ApiKey string, Pin string, Version int, Opts Options) {
	blockIo.defaultVersion = 2
	blockIo.defaultServer = ""
	blockIo.defaultPort = ""
	blockIo.host = "block.io"

	if Opts.allowNoPin == "" {
		Opts.allowNoPin = "false"
	}
	blockIo.options = Opts
	blockIo.apiUrl = blockIo.options.apiUrl

	blockIo.pin = Pin
	blockIo.aesKey = ""

	blockIo.apiKey = ApiKey
	if Version == -1 {
		blockIo.version = blockIo.defaultVersion
	} else {
		blockIo.version = Version
	}
	blockIo.server = blockIo.defaultServer
	blockIo.port = blockIo.defaultPort
	if Pin != "" {
		blockIo.pin = Pin
		blockIo.aesKey = lib.PinToAesKey(blockIo.pin)
	}

	var ServerString string
	if blockIo.server != "" {
		ServerString = blockIo.server + "."
	} else {
		ServerString = blockIo.server
	}

	var PortString string
	if blockIo.port != "" {
		PortString = ":" + blockIo.port
	} else {
		PortString = blockIo.port
	}

	if blockIo.apiUrl == "" {
		blockIo.apiUrl = "https://" + ServerString + blockIo.host + PortString + "/api/v" + strconv.Itoa(blockIo.version) + "/"
	}
	blockIo.restClient = resty.New()
}

func (blockIo *Client) get(path string) string {
	client := resty.New()

	resp, _ := client.R().
		EnableTrace().
		Get(blockIo.constructUrl(path))

	return resp.String()

}
func (blockIo *Client) post(jsonInput string, path string) string {
	var argsObj map[string]interface{}
	_ = json.Unmarshal([]byte(jsonInput), &argsObj)


	client := resty.New()

	resp, _ := client.R().
		SetHeader("Accept", "application/json").
		SetBody(argsObj).
		Post(blockIo.constructUrl(path))

	return resp.String()
}

func (blockIo *Client) _withdraw(Method string, Path string, args string) map[string]interface{} {

	var argsObj map[string]interface{}

	err := json.Unmarshal([]byte(args), &argsObj)

	if err != nil {
		panic(err)
	}

	var pin string
	if argsObj["pin"] != nil {
		pin = argsObj["pin"].(string)
	} else {
		pin = blockIo.pin
	}
	argsObj["pin"] = ""
	if pin != "" {
	}
	res := blockIo._request(Method, Path, args)

	jsonString, _ := json.Marshal(res)
	var pojo SignatureData
	err = json.Unmarshal([]byte(string(jsonString)), &pojo)
	if err != nil {
		panic(err)
	}
	if (pojo.ReferenceID == "" || pojo.EncryptedPassphrase == EncryptedPassphrase{} ||
		pojo.EncryptedPassphrase.Passphrase == "") {
		return res
	}
	if pin == "" {
		if blockIo.options.allowNoPin == "true" {
			return res
		}
	}
	var encryptedPassphrase = pojo.EncryptedPassphrase.Passphrase
	var aesKey string
	if blockIo.aesKey != "" { aesKey = blockIo.aesKey
	} else {aesKey = lib.PinToAesKey(pin) }
	privKey := lib.ExtractKeyFromEncryptedPassphrase(encryptedPassphrase,aesKey)
	pubKey := lib.ExtractPubKeyFromEncryptedPassphrase(encryptedPassphrase,aesKey)
	if pubKey != pojo.EncryptedPassphrase.SignerPublicKey {
		panic("Public key mismatch. Invalid Secret PIN detected.")
	}

	for i := 0; i < len(pojo.Inputs); i++ {
		for j := 0; j < len(pojo.Inputs[i].Signers); j++ {
			pojo.Inputs[i].Signers[j].SignedData = lib.SignInputs(privKey,pojo.Inputs[i].DataToSign)
		}
	}
	pojo.EncryptedPassphrase = EncryptedPassphrase{}
	pojoMarshalled, _ := json.Marshal(pojo)
	return blockIo._request(Method,"sign_and_finalize_withdrawal",string(pojoMarshalled))
}

func (blockIo *Client) _sweep(Method string, Path string, args string) map[string]interface{} {

	var argsObj map[string]interface{}

	err := json.Unmarshal([]byte(args), &argsObj)

	if err != nil {
		panic(err)
	}

	if argsObj["to_address"] == nil {
		panic("Missing mandatory private_key argument.")
	}

	privKeyStr := argsObj["private_key"].(string)
	keyFromWif := lib.FromWIF(privKeyStr)
	argsObj["public_key"] = lib.PubKeyFromWIF(privKeyStr)
	argsObj["private_key"] = ""

	argsObjMarshalled,_ := json.Marshal(argsObj)

	res := blockIo._request(Method, Path, string(argsObjMarshalled))

	jsonString, _ := json.Marshal(res)
	var pojo SignatureData
	err = json.Unmarshal([]byte(string(jsonString)), &pojo)
	if err != nil {
		panic(err)
	}
	if pojo.ReferenceID == "" {
		return res
	}

	for i := 0; i < len(pojo.Inputs); i++ {
		for j := 0; j < len(pojo.Inputs[i].Signers); j++ {
			pojo.Inputs[i].Signers[j].SignedData = lib.SignInputs(keyFromWif,pojo.Inputs[i].DataToSign)
		}
	}
	pojo.EncryptedPassphrase = EncryptedPassphrase{}
	pojoMarshalled, _ := json.Marshal(pojo)
	return blockIo._request(Method,"sign_and_finalize_sweep",string(pojoMarshalled))
}

func (blockIo *Client) _request(Method string, Path string, args string) map[string]interface{} {
	var res map[string]interface{}
	if Method == "POST" {
		if strings.Contains(Path, "sign_and_finalize") {

			var postObj map[string]string = map[string]string{"signature_data": args}

			temp, _ := json.Marshal(postObj)
			args = string(temp)
		}
		resString := blockIo.post(args, Path)
	_:
		json.Unmarshal([]byte(resString), &res)
	} else {
		resString := blockIo.get(Path)
	_:
		json.Unmarshal([]byte(resString), &res)
	}
	jsonString, _ := json.Marshal(res["data"])
	res = nil
	err := json.Unmarshal([]byte(string(jsonString)), &res)
	if err != nil {
		panic(err)
	}
	if res == nil {
		panic("No response from API server")
	}
	return res
}

func (blockIo *Client) constructUrl(path string) string {
	return blockIo.apiUrl + path + "?api_key=" + blockIo.apiKey
}

