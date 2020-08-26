package block_io_go

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/go-resty/resty/v2"
	"github.com/piotrnar/gocoin/lib/btc"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/exp/utf8string"
	"log"
)

func SignInputs(PrivKey string, DataToSign string) string {
	// Decode a hex-encoded private key.
	pkBytes, err := hex.DecodeString(PrivKey)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)
	pubKey = pubKey
	// Sign a message using the private key.
	message := DataToSign
	messageHash, _ := hex.DecodeString(message)
	signature, err := privKey.Sign(messageHash)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return hex.EncodeToString(signature.Serialize())
}

func ParseResponseData(res *resty.Response) (SignatureData, error){
	var withdrawRes SignatureRes
	marshalErr := json.Unmarshal([]byte(res.String()), &withdrawRes)
	if marshalErr != nil {
		return SignatureData{}, marshalErr
	}
	return withdrawRes.Data, nil
}

func ExtractKeyFromEncryptedPassphrase(EncryptedData string, B64Key string) string {
	Decrypted := Decrypt(EncryptedData,B64Key)
	Unhexlified, err := hex.DecodeString(Decrypted)

	if err != nil {
		log.Fatal(errors.New("Unhexlified Error"))
	}

	Hashed := sha256.Sum256(Unhexlified)
	UsableHashed := Hashed[:]
	return hex.EncodeToString(UsableHashed)
}

func ExtractPubKeyFromEncryptedPassphrase(EncryptedData string, B64Key string) string {
	Decrypted := Decrypt(EncryptedData,B64Key)
	Unhexlified, err := hex.DecodeString(Decrypted)

	if err != nil {
		log.Fatal(errors.New("Unhexlified Error"))
	}

	Hashed := sha256.Sum256(Unhexlified)
	UsableHashed := Hashed[:]

	result := btc.PublicFromPrivate(UsableHashed, true)

	return hex.EncodeToString(result)
}

func PinToAesKey(pin string) string {
	var saltOld []byte = make([]byte, 0)
	var salt [1024]byte;
	for i := 0; i < 1024; i++ {
		salt[i] = 0
	}
	pinBytes := []byte(pin)

	firstHash := hex.EncodeToString(pbkdf2.Key(pinBytes, saltOld, 1024, 16, sha256.New))

	firstHashBytes := []byte(firstHash)

	key := pbkdf2.Key(firstHashBytes, saltOld, 1024, 32, sha256.New)

	return base64.StdEncoding.EncodeToString(key)
}

func SHA256_hash(ba []byte) []byte {
	sha := sha256.Sum256(ba)
	return sha[:]
}

type ecb struct {
	b         cipher.Block
	blockSize int
}

type ecbEncrypter ecb

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}

func (x *ecbEncrypter) BlockSize() int { return x.blockSize }
func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}
func pKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func AESEncrypt(src string, key string) []byte {
	usableKey, _ := base64.StdEncoding.DecodeString(key)
	block, err := aes.NewCipher(usableKey)
	if err != nil {
		fmt.Println("key error1", err)
	}
	if src == "" {
		fmt.Println("plain content empty")
	}
	ecb := NewECBEncrypter(block)
	content := []byte(utf8string.NewString(src).String())

	content = pKCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))

	ecb.CryptBlocks(crypted, content)

	return crypted
}

type ecbDecrypter ecb

func (x *ecbDecrypter) BlockSize() int { return x.blockSize }
func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}

func pKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func AESDecrypt(crypt string, key string) string {
	usableKey, _ := base64.StdEncoding.DecodeString(key)

	block, err := aes.NewCipher(usableKey)
	if err != nil {
		fmt.Println("key error1", err)
	}
	if len(crypt) == 0 {
		fmt.Println("plain content empty")
	}
	usableCrypt, _ := base64.StdEncoding.DecodeString(crypt)

	ecb := NewECBDecrypter(block)
	decrypted := make([]byte, len(usableCrypt))
	ecb.CryptBlocks(decrypted, []byte(usableCrypt))

	padded := pKCS5UnPadding(decrypted)

	return utf8string.NewString(string(padded)).String()
}

func Encrypt(data string, key string) string {

	temp := AESEncrypt(data, key)
	return base64.StdEncoding.EncodeToString(temp)
}

func Decrypt(data string, key string) string {
	temp := AESDecrypt(data, key)
	return temp
}
