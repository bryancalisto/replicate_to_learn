package main

import (
	"encoding/hex"
	"fmt"
)

func main() {
	data := "hello"
	result := hash([]byte(data))
	fmt.Printf("DIGEST: %s\n", hex.EncodeToString(result[:]))
}
