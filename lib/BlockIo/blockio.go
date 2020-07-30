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
	Options        map[string]interface{}
	ApiUrl         string
	Pin            string
	AesKey         string
	ApiKey         string
	Version        int
	DefaultVersion int
	Server         string
	DefaultServer  string
	Port           string
	DefaultPort    string
	Host           string
	RestClient     *resty.Client
}

func (blockIo *Client) Instantiate(Config string, Pin string, Version int, Options string) {
	blockIo.DefaultVersion = 2
	blockIo.DefaultServer = ""
	blockIo.DefaultPort = ""
	blockIo.Host = "block.io"

	if Options == "" {
		Options = "{}"
	}

	_ = json.Unmarshal([]byte(`{"allowNoPin": false}`), &blockIo.Options)
	blockIo.ApiUrl = ""
	blockIo.Pin = Pin
	blockIo.AesKey = ""
	var ConfigObj map[string]interface{}

	//no err returned, meaning valid json string in Config
	if json.Unmarshal([]byte(Config), &ConfigObj) == nil {
		blockIo.ApiKey = ConfigObj["api_key"].(string)
		if ConfigObj["version"] != nil {
			blockIo.Version = ConfigObj["version"].(int)
		} else {
			blockIo.Version = blockIo.DefaultVersion
		}
		if ConfigObj["server"] != nil {
			blockIo.Server = ConfigObj["server"].(string)
		} else {
			blockIo.Server = blockIo.DefaultServer
		}
		if ConfigObj["port"] != nil {
			blockIo.Port = ConfigObj["port"].(string)
		} else {
			blockIo.Port = blockIo.DefaultPort
		}

		if ConfigObj["pin"] != nil {
			blockIo.Pin = ConfigObj["pin"].(string)
			blockIo.AesKey = lib.PinToAesKey(blockIo.Pin)
		}
		if ConfigObj["options"] != nil {
			var Options map[string]interface{}
			err := json.Unmarshal([]byte(ConfigObj["options"].(string)), &Options)

			if err != nil {
				fmt.Println("Ignoring invalid options.", err)
			}
			blockIo.Options = Options
			blockIo.Options["allowNoPin"] = false
		}

	} else {
		//else block will be accessed if Config is not a valid json string
		blockIo.ApiKey = Config
		if Version == -1 {
			blockIo.Version = blockIo.DefaultVersion
		} else {
			blockIo.Version = Version
		}
		blockIo.Server = blockIo.DefaultServer
		blockIo.Port = blockIo.DefaultPort
		if Pin != "" {
			blockIo.Pin = Pin
			blockIo.AesKey = lib.PinToAesKey(blockIo.Pin)
		}
	}

	if Options != "" {
		err := json.Unmarshal([]byte(Options), &blockIo.Options)
		if err != nil {
			fmt.Println("Ingoring invalid options", err)
		}
		blockIo.Options["allowNoPin"] = false
		if blockIo.Options["api_url"] != nil {
			blockIo.ApiUrl = blockIo.Options["api_url"].(string)
		} else {
			blockIo.ApiUrl = ""
		}
	}

	var ServerString string
	if blockIo.Server != "" {
		ServerString = blockIo.Server + "."
	} else {
		ServerString = blockIo.Server
	}

	var PortString string
	if blockIo.Port != "" {
		PortString = ":" + blockIo.Port
	} else {
		PortString = blockIo.Port
	}

	if blockIo.ApiUrl == "" {
		blockIo.ApiUrl = "https://" + ServerString + blockIo.Host + PortString + "/api/v" + strconv.Itoa(Version) + "/"
	}

	blockIo.RestClient = resty.New()

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
		pin = blockIo.Pin
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
		 if blockIo.Options["allowNoPin"] == true {
			return res
		 }
	}
	var encrypted_passphrase string = pojo.EncryptedPassphrase.Passphrase
	var aesKey string
	if blockIo.AesKey != "" { aesKey = blockIo.AesKey } else {aesKey = lib.PinToAesKey(pin) }
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

	/**
	TODO: define sweep here
	 */
	return map[string]interface{}{}

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
	return blockIo.ApiUrl + path + "?api_key=" + blockIo.ApiKey
}
