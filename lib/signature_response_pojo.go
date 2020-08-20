package lib

type SignatureData struct {
	EncryptedPassphrase  EncryptedPassphrase `json:"encrypted_passphrase"`
	Inputs               Inputs              `json:"inputs"`
	MoreSignaturesNeeded bool                `json:"more_signatures_needed"`
	ReferenceID          string              `json:"reference_id"`
	UnsignedTxHex        string              `json:"unsigned_tx_hex"`
}

type EncryptedPassphrase struct {
	Passphrase      string `json:"passphrase"`
	SignerAddress   string `json:"signer_address"`
	SignerPublicKey string `json:"signer_public_key"`
}

type Inputs []struct {
	DataToSign       string  `json:"data_to_sign"`
	InputNo          int64   `json:"input_no"`
	SignaturesNeeded int64   `json:"signatures_needed"`
	Signers          Signers `json:"signers"`
}

type Signers []struct {
	SignedData      interface{} `json:"signed_data"`
	SignerAddress   string      `json:"signer_address"`
	SignerPublicKey string      `json:"signer_public_key"`
}