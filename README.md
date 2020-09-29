# BlockIo

This Golang library is the official Block.IO low-level SDK. To use the functions
provided by this SDK, you will need a REST client of your choice, the
Bitcoin, Litecoin, or Dogecoin API key(s) from <a href="https://block.io" target="_blank">Block.io</a>
and, if required by your use case, your PIN or secret keys.

## Installation

```bash
  go get github.com/BlockIo/block_io-go
```

## Usage

### Example implementations

- [Create new address](https://github.com/BlockIo/block_io-go/tree/master/examples/get_address)
- [Get balance](https://github.com/BlockIo/block_io-go/tree/master/examples/get_balance)
- [Withdraw](https://github.com/BlockIo/block_io-go/tree/master/examples/withdraw)
- [Sweep external wallet](https://github.com/BlockIo/block_io-go/tree/master/examples/sweep)
- [DTrust integration](https://github.com/BlockIo/block_io-go/tree/master/examples/dtrust)

## Method documentation

### Signing

#### SignWithdrawRequestJson()

```go
  func SignWithdrawRequestJson(pin string, withdrawData string) (string, error)
```

Signs JSON encoded signature requests returned from the `/api/v2/withdraw*`
endpoints with a PIN-derived key and returns a JSON encoded string that can be
posted to `/api/v2/sign_and_finalize_withdrawal`.

#### SignRequestJsonWithKey()

```go
  func SignRequestJsonWithKey(ecKey *ECKey, data string) (string, error)
```

Signs JSON encoded strings returned from the `/api/v2/sweep*` and
`/api/v2/withdraw_dtrust*/` endpoints with a local ECKey and returns a JSON
encoded string that can be posted to `/api/v2/sign_and_finalize_*` endpoints.

#### SignRequestJsonWithKeys()

```go
  func SignRequestJsonWithKeys(ecKeys []*ECKey, data string) (string, error)
```

Signs JSON encoded strings returned from the `/api/v2/withdraw_dtrust*/`
endpoints with multiple local ECKeys and returns a JSON encoded string that can
be posted to `/api/v2/sign_and_finalize_*` endpoints.

### Key derivation

#### NewECKey()

```go
  func NewECKey (d [32]byte, compressed bool) *ECKey
```

Creates an ECKey from a byte slice.

#### FromWIF()

```go
  func FromWIF(strWif string) (*ECKey, error)
```

Creates an ECKey from a WIF-formatted string. Returns an error if the given
string could not be

#### DeriveKeyFromHex()

```go
  func DeriveKeyFromHex(hexPass string) (*ECKey, error)
```

Derives an ECKey from a hexadecimal encoded string by hashing it with sha256 once.

#### DeriveKeyFromString()

```go
  func DeriveKeyFromString(pass string) *ECKey
```

Convenience function to derive an ECKey from a string seed by hashing it with
sha256 once.

_NOTE: Do not use this for deriving production keys!_

## Testing

```bash
  go test -v
```
