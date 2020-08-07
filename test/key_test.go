package test

import (
	"github.com/BlockIo/block_io-go/lib"
	"testing"
)

var wif string
var passphrase string
var dataToSign string
var controlPrivKeyFromPassphrase string
var controlPubKeyFromPassphrase string
var controlPrivKeyFromWif string
var controlPubKeyFromWif string
var controlSignedDataWifKey string
var controlSignedDataPassphraseKey string

var privKeyFromWif string
var pubKeyFromWif string
var privKeyFromPassphrase string
var pubKeyFromPassphrase string
var signedDataWifKey string
var signedDataPassphraseKey string

func Setup() {
	wif = "L1cq4uDmSKMiViT4DuR8jqJv8AiiSZ9VeJr82yau5nfVQYaAgDdr";
	passphrase = "deadbeef";
	dataToSign = "e76f0f78b7e7474f04cc14ad1343e4cc28f450399a79457d1240511a054afd63";
	controlPrivKeyFromPassphrase = "5f78c33274e43fa9de5659265c1d917e25c03722dcb0b8d27db8d5feaa813953";
	controlPubKeyFromPassphrase = "02953b9dfcec241eec348c12b1db813d3cd5ec9d93923c04d2fa3832208b8c0f84";
	controlPrivKeyFromWif = "833e2256c42b4a41ee0a6ee284c39cf8e1978bc8e878eb7ae87803e22d48caa9";
	controlPubKeyFromWif = "024988bae7e0ade83cb1b6eb0fd81e6161f6657ad5dd91d216fbeab22aea3b61a0";
	controlSignedDataWifKey = "3045022100aec97f7ad7a9831d583ca157284a68706a6ac4e76d6c9ee33adce6227a40e675022008894fb35020792c01443d399d33ffceb72ac1d410b6dcb9e31dcc71e6c49e92";
	controlSignedDataPassphraseKey = "30450221009a68321e071c94e25484e26435639f00d23ef3fbe9c529c3347dc061f562530c0220134d3159098950b81b678f9e3b15e100f5478bb45345d3243df41ae616e70032";
	privKeyFromWif, _ = lib.FromWIF(wif)
	pubKeyFromWif = lib.PubKeyFromWIF(wif)
	privKeyFromPassphrase = lib.ExtractKeyFromPassphrase(passphrase)
	pubKeyFromPassphrase = lib.ExtractPubKeyFromPassphrase(passphrase)
	signedDataWifKey = lib.SignInputs(privKeyFromWif, dataToSign)
	signedDataPassphraseKey = lib.SignInputs(privKeyFromPassphrase, dataToSign)
}

func TestPrivKeyFromWif(t *testing.T) {
	Setup()
	if privKeyFromWif != controlPrivKeyFromWif {
		t.Error("fromWIF not giving correct value")
	}
}

func TestPubKeyFromWif(t *testing.T) {
	if pubKeyFromWif != controlPubKeyFromWif {
		t.Error("public key from wif not giving correct value")
	}
}

func TestPubKeyFromPassphrase(t *testing.T) {
	if pubKeyFromPassphrase != controlPubKeyFromPassphrase {
		t.Error("public key from passphrase not giving correct value")
	}
}

func TestPrivKeyFromPassphrase(t *testing.T) {
	if privKeyFromPassphrase != controlPrivKeyFromPassphrase {
		t.Error("from passphrase not giving correct value")
	}
}

func TestSignDataWifKey(t *testing.T) {
	if signedDataWifKey != controlSignedDataWifKey {
		t.Error("signed data from wif key not giving correct value")
	}
}

func TestSignDataPassphraseKey(t *testing.T) {
	if signedDataPassphraseKey != controlSignedDataPassphraseKey {
		t.Error("signed data from passphrase key not giving correct value")
	}
}
