package block_io_go

import (
	"testing"
	"encoding/base64"
)

var pin string
var aesKey string
var controlClearText string
var controlCipherText string

func HelperSetup() {
	pin = "123456";
	aesKey = PinToAesKey(pin);
	controlClearText = "I'm a little tea pot short and stout";
	controlCipherText = "7HTfNBYJjq09+vi8hTQhy6lCp3IHv5rztNnKCJ5RB7cSL+NjHrFVv1jl7qkxJsOg";
}

func TestPinToAes(t *testing.T) {
	var controlData = "0EeMOVtm5YihUYzdCNgleqIUWkwgvNBcRmr7M0t9GOc=";
	HelperSetup()
	if aesKey != controlData {
		t.Error("PinToAes not returning correct output")
	}
}

func TestEncrypt(t *testing.T) {
	usableKey, _ := base64.StdEncoding.DecodeString(aesKey)
	//TODO test for error
	encryptedData, _ := AESEncrypt([]byte(controlClearText), usableKey);
	base64Data := base64.StdEncoding.EncodeToString(encryptedData)
	if base64Data != controlCipherText {
		t.Error("Encrypt not returning correct output")
	}
}

func TestDecrypt(t *testing.T) {
	usableKey, _ := base64.StdEncoding.DecodeString(aesKey)
	usableCrypt, _ := base64.StdEncoding.DecodeString(controlCipherText)
	//TODO test for error
	decryptedData, _ := AESDecrypt(usableCrypt, usableKey)
	string := string(decryptedData)
	if string != controlClearText {
		t.Error("Decrypt not returning correct output")
	}
}
