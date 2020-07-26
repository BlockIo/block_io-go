package main

import (
	"fmt"
	"github.com/BlockIo/block_io-go/lib"
)

func main() {
	test := lib.ByteArrayToHexString([]byte {0})
	fmt.Println("BlockIo Lib: " + test)
}