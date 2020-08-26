package block_io_go

import (
	"encoding/hex"
	"testing"
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

func TestSha256(t *testing.T) {
	HelperSetup()
	controlData := "5f78c33274e43fa9de5659265c1d917e25c03722dcb0b8d27db8d5feaa813953";
	testData := "deadbeef";
	bytes, _ := hex.DecodeString(testData);
	shaData := hex.EncodeToString(SHA256_hash(bytes));
	if shaData != controlData {
		t.Error("SHA256 not returning correct output")
	}
}

func TestPinToAes(t *testing.T) {
	var controlData = "0EeMOVtm5YihUYzdCNgleqIUWkwgvNBcRmr7M0t9GOc=";
	if aesKey != controlData {
		t.Error("PinToAes not returning correct output")
	}
}

func TestEncrypt(t *testing.T) {
	encryptedData := Encrypt(controlClearText, aesKey);
	if encryptedData != controlCipherText {
		t.Error("Encrypt not returning correct output")
	}
}

func TestDecrypt(t *testing.T) {
	decryptedData := Decrypt(controlCipherText, aesKey)
	string := string(decryptedData)
	if string != controlClearText {
		t.Error("Decrypt not returning correct output")
	}
}
