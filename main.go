package main

import (
	"ChaCha20/internal"
	"fmt"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(6)
	key := internal.Key{1, 2, 3, 4, 5, 6, 7, 8}
	nonce := internal.Nonce{1, 2, 3}

	str := "aljsdladklakdakldasjdanidnasdalkjhdahdiapdjawpdmaiudnipmasdhashdpajsdnasčodhasdošpjasdouabdpiasjndšasodkjaiphdaspdasšpdokšasdmapsdniasnšoasmcšađadadaćsčda\nasdasd"
	data := []byte(str)
	fmt.Println("Original str: " + str)

	e := internal.Encrypt(key, 0, nonce, data)
	fmt.Println("Encrypted string: " + string(e))

	e2 := []byte(e)
	ed := internal.Encrypt(key, 0, nonce, e2)
	fmt.Println("Decrypted string: " + string(ed))
}
