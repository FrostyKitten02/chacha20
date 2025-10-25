package main

import (
	"ChaCha20/internal"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	runtime.GOMAXPROCS(6)

	key := internal.Key{1, 2, 3, 4, 5, 6, 7, 8}
	nonce := internal.Nonce{1, 2, 3}

	file, _ := os.ReadFile("test.mkv")
	encrypt(key, 0, nonce, file, internal.EncryptParallel, "parallel")
	//os.WriteFile("test.mkv.enc", out, 0777)
	//
	//back := encrypt(key, 0, nonce, out, internal.EncryptParallel, "parallel-dec")
	//os.WriteFile("test2.mkv", back, 0777)

	//encrypt(key, 0, nonce, file, internal.Encrypt, "sequential")

}

func encrypt(key internal.Key, counter uint32, nonce internal.Nonce, data []byte, e func(internal.Key, uint32, internal.Nonce, []byte) []byte, name string) []byte {
	size := len(data)
	megaSize := float64(size) / 1_000_000
	start := time.Now().UnixNano()
	res := e(key, counter, nonce, data)
	totalTime := float64(time.Now().UnixNano()-start) / 1_000_000_000
	fmt.Println("Encrypted " + name + " in: " + strconv.FormatFloat(totalTime, 'f', 3, 64) + "s")
	fmt.Println("Encrypt " + name + ": " + strconv.FormatFloat(megaSize/totalTime, 'f', 3, 64) + "MB/s")

	return res
}
