# block_io-go

# BlockIo

This Golang library is the official Block.IO SDK. To call the API, you will need the Dogecoin, Bitcoin, or Litecoin API key(s) from <a href="https://block.io" target="_blank">Block.io</a>. Go ahead, sign up :)

## Installation

    go get github.com/BlockIo/block_io-go

## Usage

It's super easy to get started. In your code, do this:

    import blockio "github.com/BlockIo/block_io-go"

    // Withdraw json response signing

    withdrawData, _ := blockio.ParseResponseData(rawWithdrawResponse.String())
	signatureReq, _ := blockio.SignWithdrawRequest(pin, withdrawData)

    // Sweep json response signing

    ecKey, _ := blockio.FromWIF(strWif)
    sweepData, _ := blockio.ParseResponseData(rawSweepResponse.String())
	signatureReq, _ := blockio.SignSweepRequest(ecKey, sweepData)

##### For a more detailed guide on usage, check the examples folder in the repo 

## Testing

```bash
go test -v
