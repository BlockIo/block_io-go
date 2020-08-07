# block_io-go

# BlockIo

This Golang library is the official reference client for the Block.io payments API and uses go modules. To use this, you will need the Dogecoin, Bitcoin, or Litecoin API key(s) from <a href="https://block.io" target="_blank">Block.io</a>. Go ahead, sign up :)

## Installation

1. Clone the repo
2. go get -v ./lib

## Usage

It's super easy to get started. In your code, do this:

    var blockIo BlockIo.Client
    blockIo.Instantiate(apiKey, pin, version, BlockIo.Options{})

    // print the account balance request's response
    res, err := blockIo.GetBalance(map[string]interface{}{})

    // print the response of a withdrawal request
    // 'SECRET_PIN' is only required if you did not specify it at 
    // class initialization time.
    res, _ := blockIo.Withdraw(map[string]interface{}{
		"from_labels": "default",
		"to_label": "testDest15",
		"amount":"2.5",
        "pin":pin,
	})  

##### For a more detailed guide on usage, check the examples folder in the repo 

##### A note on passing json args to requests:

Args are passed as map[string]interface{} 

    map[string]interface{}{ "param": "val", "array": "val1, val2" }

## Testing

**DO NOT USE PRODUCTION CREDENTIALS FOR UNIT TESTING!** 

Test syntax:

```bash
go test -v ./test
