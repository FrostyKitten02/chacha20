package internal

import (
	"sync"
)

func EncryptParallel(key Key, counter uint32, nonce Nonce, data []byte) []byte {
	dataLen := uint32(len(data))
	maxBlock := dataLen / 64
	encrypted := make([]byte, dataLen)

	numWorkers := 10
	jobs := make(chan uint32, numWorkers*2)
	var wg sync.WaitGroup

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go worker(key, counter, nonce, data, dataLen, encrypted, jobs, &wg)
	}

	for i := uint32(0); i <= maxBlock; i++ {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
	return encrypted
}

func worker(key Key, counter uint32, nonce Nonce, data []byte, dataLen uint32, encrypted []byte, jobs <-chan uint32, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := range jobs {
		key_stream := blockFunc(key, counter+i, nonce)
		startIndex := i * 64
		finishIndex := startIndex + 64

		if finishIndex > dataLen {
			finishIndex = dataLen
		}

		block := data[startIndex:finishIndex]
		xorToArr(key_stream, block, encrypted, int(startIndex))
	}
}
