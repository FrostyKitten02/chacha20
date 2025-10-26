package chacha20

import (
	"ChaCha20/internal"
	"io"
	"os"
	"sync"
)

const chunkSize = 320 * 1024 * 1024

func EncryptFile(key internal.Key, counter uint32, nonce internal.Nonce, inputPath, outputPath string) error {
	inFile, inFileErr := os.Open(inputPath)
	if inFileErr != nil {
		return inFileErr
	}
	defer inFile.Close()

	outFile, outFileErr := os.Create(outputPath)
	if outFileErr != nil {
		return outFileErr
	}
	defer outFile.Close()

	info, statErr := inFile.Stat()
	if statErr != nil {
		return statErr
	}
	preallocateFileErr := outFile.Truncate(info.Size())
	if preallocateFileErr != nil {
		return preallocateFileErr
	}

	buffer := make([]byte, chunkSize)
	currentCounter := counter

	for {
		n, err := inFile.Read(buffer)
		if n > 0 {
			chunk := buffer[:n]
			encrypted := EncryptParallel(key, currentCounter, nonce, chunk)

			_, writeErr := outFile.Write(encrypted)
			if writeErr != nil {
				return writeErr
			}

			blocksInChunk := uint32((n + 63) / 64)
			currentCounter += blocksInChunk
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func EncryptParallel(key internal.Key, counter uint32, nonce internal.Nonce, data []byte) []byte {
	dataLen := len(data)
	maxBlock := dataLen / 64
	encrypted := make([]byte, dataLen)

	numWorkers := 10
	jobs := make(chan int, numWorkers*2)
	var wg sync.WaitGroup

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go internal.Worker(key, counter, nonce, data, dataLen, encrypted, jobs, &wg)
	}

	for i := 0; i <= maxBlock; i++ {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
	return encrypted
}

func Encrypt(key internal.Key, counter uint32, nonce internal.Nonce, data []byte) []byte {
	dataLen := len(data)
	maxBlock := dataLen / 64
	encrypted := make([]byte, dataLen)
	for i := 0; i <= maxBlock; i++ {
		internal.EncryptBlock(key, counter, nonce, data, i, dataLen, encrypted)
	}

	return encrypted
}
