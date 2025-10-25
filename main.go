package main

import (
	"ChaCha20/internal"
	"fmt"
)

func main() {
	key := internal.Key{1, 2, 3, 4, 5, 6, 7, 8}
	nonce := internal.Nonce{1, 2, 3}

	data := []byte{123, 10, 50, 12}

	e := internal.Encrypt(key, 0, nonce, data)
	ed := internal.Encrypt(key, 0, nonce, e)

	fmt.Println(e)
	fmt.Println(ed)

}
