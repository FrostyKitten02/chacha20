package internal

import (
	"sync"
)

func EncryptParallel(key Key, counter uint32, nonce Nonce, data []byte) []byte {
	dataLen := len(data)
	maxBlock := dataLen / 64
	encrypted := make([]byte, dataLen)

	numWorkers := 10
	jobs := make(chan int, numWorkers*2)
	var wg sync.WaitGroup

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go worker(key, counter, nonce, data, dataLen, encrypted, jobs, &wg)
	}

	for i := 0; i <= maxBlock; i++ {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
	return encrypted
}

func worker(key Key, counter uint32, nonce Nonce, data []byte, dataLen int, encrypted []byte, jobs <-chan int, wg *sync.WaitGroup) {
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
		xorToArr(key_stream, block, encrypted, int(startIndex))
	}
}
