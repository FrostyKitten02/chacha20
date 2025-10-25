package internal

import (
	"io"
	"os"
)

const chunkSize = 320 * 1024 * 1024

func EncryptFile(key Key, counter uint32, nonce Nonce, inputPath, outputPath string) error {
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
