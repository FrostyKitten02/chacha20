package internal

import (
	"sync"
)

func Worker(key Key, counter uint32, nonce Nonce, data []byte, dataLen int, encrypted []byte, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := range jobs {
		//same code as in single thread function but not calling that function because of 20% overhead for some reason!
		key_stream := blockFunc(key, uint32(int(counter)+i), nonce)
		startIndex := int(i) * 64
		finishIndex := startIndex + 64

		if finishIndex > dataLen {
			finishIndex = dataLen
		}

		block := data[startIndex:finishIndex]
		xorToArr(key_stream, block, encrypted, startIndex)
	}
}
