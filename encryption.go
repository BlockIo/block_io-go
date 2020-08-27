package block_io_go

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

type ecb struct {
	cipher         cipher.Block
}

func newECB(c cipher.Block) *ecb {
	return &ecb{
		cipher:         c,
	}
}

func (x *ecb) BlockSize() int {
	return x.cipher.BlockSize()
}

func (x *ecb) EncryptBlocks(dst []byte, src []byte) error {

	if len(src) % x.BlockSize() != 0 {
		return errors.New("EncryptBlocks: input not a multiple of blocksize")
	}

	if len(dst) < len(src) {
		return errors.New("EncryptBlocks: output buffer is too small")
	}

	for len(src) > 0 {
		x.cipher.Encrypt(dst, src[:x.BlockSize()])
		src = src[x.BlockSize():]
		dst = dst[x.BlockSize():]
	}

	return nil
}

func (x *ecb) DecryptBlocks(dst []byte, src []byte) error {

	if len(src) % x.BlockSize() != 0 {
		return errors.New("DecryptBlocks: input not a multiple of blocksize")
	}

	if len(dst) < len(src) {
		return errors.New("DecryptBlocks: output buffer is too small")
	}

	for len(src) > 0 {
		x.cipher.Decrypt(dst, src[:x.BlockSize()])
		src = src[x.BlockSize():]
		dst = dst[x.BlockSize():]
	}

	return nil
}

func (x *ecb) pkcs5Padding(ciphertext []byte) []byte {
	padding := x.BlockSize() - len(ciphertext) % x.BlockSize()
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func (x *ecb) pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func AESEncrypt(clearText []byte, key []byte) ([]byte, error) {

	cipher, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	if len(clearText) == 0 {
		return nil, errors.New("AESEncrypt: Cleartext length is 0")
	}

	ecb := newECB(cipher)

	padded := ecb.pkcs5Padding(clearText)
	cipherText := make([]byte, len(padded))

	ecb.EncryptBlocks(cipherText, padded)

	return cipherText, nil
}

func AESDecrypt(cipherText []byte, key []byte) ([]byte, error) {
	cipher, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	if len(cipherText) == 0 {
		return nil, errors.New("AESDecrypt: Ciphertext length is 0")
	}

	ecb := newECB(cipher)
	decrypted := make([]byte, len(cipherText))
	ecb.DecryptBlocks(decrypted, []byte(cipherText))

	unpadded := ecb.pkcs5UnPadding(decrypted)

	return unpadded, nil
}
