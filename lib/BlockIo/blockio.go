package BlockIo

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
	"strings"

	"github.com/BlockIo/block_io-go/lib"
)

type Client struct {
	options        map[string]interface{}
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

func (blockIo *Client) Instantiate(Config string, Pin string, Version int, Options string) {
	blockIo.defaultVersion = 2
	blockIo.defaultServer = ""
	blockIo.defaultPort = ""
	blockIo.host = "block.io"

	if Options == "" {
		Options = "{}"
	}

	_ = json.Unmarshal([]byte(`{"allowNoPin": false}`), &blockIo.options)
	blockIo.apiUrl = ""
	blockIo.pin = Pin
	blockIo.aesKey = ""
	var ConfigObj map[string]interface{}

	err := json.Unmarshal([]byte(Config), &ConfigObj)
	//no err returned, meaning valid json string in Config
	if err == nil {
		blockIo.apiKey = ConfigObj["api_key"].(string)
		if ConfigObj["version"] != nil {
			blockIo.version = int(ConfigObj["version"].(float64))
		} else {
			blockIo.version = blockIo.defaultVersion
		}
		if ConfigObj["server"] != nil {
			blockIo.server = ConfigObj["server"].(string)
		} else {
			blockIo.server = blockIo.defaultServer
		}
		if ConfigObj["port"] != nil {
			blockIo.port = ConfigObj["port"].(string)
		} else {
			blockIo.port = blockIo.defaultPort
		}

		if ConfigObj["pin"] != nil {
			blockIo.pin = ConfigObj["pin"].(string)
			blockIo.aesKey = lib.PinToAesKey(blockIo.pin)
		}
		if ConfigObj["options"] != nil {
			var Options map[string]interface{}
			stringMap,_ := json.Marshal(ConfigObj["options"])
			_ = json.Unmarshal([]byte(string(stringMap)),&Options)

			blockIo.options = Options
			blockIo.options["allowNoPin"] = false
		}

	} else {
		//else block will be accessed if Config is not a valid json string
		blockIo.apiKey = Config
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
	}

	if Options != "" {
		err := json.Unmarshal([]byte(Options), &blockIo.options)
		if err != nil {
			fmt.Println("Ingoring invalid options", err)
		}
		blockIo.options["allowNoPin"] = false
		if blockIo.options["api_url"] != nil {
			blockIo.apiUrl = blockIo.options["api_url"].(string)
		} else {
			blockIo.apiUrl = ""
		}
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
		blockIo.apiUrl = "https://" + ServerString + blockIo.host + PortString + "/api/v" + strconv.Itoa(Version) + "/"
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
		fmt.Println("error converting json:", err)
		return nil
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

	jsonString, _ := json.Marshal(res["data"])
	var pojo SignatureData
	err = json.Unmarshal([]byte(string(jsonString)), &pojo)
	if err != nil {
		fmt.Println("Invalid response conversion:", err)
		return nil
	}
	if (pojo.ReferenceID == "" || pojo.EncryptedPassphrase == EncryptedPassphrase{} ||
		pojo.EncryptedPassphrase.Passphrase == "") {
		return res
	}
	if pin == "" {
		if blockIo.options["allowNoPin"] == true {
			return res
		}
	}
	var encrypted_passphrase string = pojo.EncryptedPassphrase.Passphrase
	var aesKey string
	if blockIo.aesKey != "" { aesKey = blockIo.aesKey
	} else {aesKey = lib.PinToAesKey(pin) }
	privKey := lib.ExtractKeyFromEncryptedPassphrase(encrypted_passphrase,aesKey)
	pubKey := lib.ExtractPubKeyFromEncryptedPassphrase(encrypted_passphrase,aesKey)
	if pubKey != pojo.EncryptedPassphrase.SignerPublicKey {
		fmt.Println("Public key mismatch. Invalid Secret PIN detected.")
		return nil
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
		fmt.Println("error converting json:", err)
		return nil
	}

	if argsObj["to_address"] == nil {
		fmt.Println("Missing mandatory private_key argument.")
	}

	privKeyStr := argsObj["private_key"].(string)
	keyFromWif := lib.FromWIF(privKeyStr)
	argsObj["public_key"] = lib.PubKeyFromWIF(privKeyStr)
	argsObj["private_key"] = ""

	argsObjMarshalled,_ := json.Marshal(argsObj)

	res := blockIo._request(Method, Path, string(argsObjMarshalled))

	jsonString, _ := json.Marshal(res["data"])
	var pojo SignatureData
	err = json.Unmarshal([]byte(string(jsonString)), &pojo)
	if err != nil {
		fmt.Println("Invalid response conversion:", err)
		return nil
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
	if res["status"] != "success" || strings.Contains(Path, "sign_and_finalize"){
		jsonString, _ := json.Marshal(res["data"])
		res = nil
		err := json.Unmarshal([]byte(string(jsonString)), &res)
		if err != nil {
			fmt.Println("Error unmarshalling response:", err)
			return nil
		}
		return res
	}

	return res
}


func (blockIo *Client) constructUrl(path string) string {
	return blockIo.apiUrl + path + "?api_key=" + blockIo.apiKey
}

