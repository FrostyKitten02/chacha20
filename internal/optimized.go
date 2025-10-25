package internal

import (
	"sync"
)

func EncryptParallel(key Key, counter uint32, nonce Nonce, data []byte) []byte {
	dataLen := uint32(len(data))
	maxBlock := dataLen / 64
	encrypted := make([]byte, dataLen)

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10)
	for i := uint32(0); i <= maxBlock; i++ {
		wg.Add(1)
		semaphore <- struct{}{}
		go encryptBlockParallel(key, counter, nonce, data, i, dataLen, encrypted, &wg, semaphore)
	}

	wg.Wait()
	return encrypted
}

func encryptBlockParallel(key Key, counter uint32, nonce Nonce, data []byte, i uint32, dataLen uint32, encrypted []byte, wg *sync.WaitGroup, semaphore chan struct{}) {
	defer wg.Done()
	defer func() { <-semaphore }()

	//encryptBlock(key, counter, nonce, data, i, dataLen, encrypted)
	key_stream := blockFunc(key, counter+i, nonce)
	startIndex := i * 64
	finishIndex := startIndex + 64

	if finishIndex > dataLen {
		finishIndex = dataLen
	}

	block := data[startIndex:finishIndex]
	xorToArr(key_stream, block, encrypted, int(startIndex))
	encryptedBlock := xorArr(key_stream, block)
	copy(encrypted[startIndex:finishIndex], encryptedBlock)
}
