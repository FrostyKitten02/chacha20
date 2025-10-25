package internal

import (
	"encoding/binary"
	"math/bits"
)

// All should be 32 bit litle endian
type Key [8]uint32
type Nonce [3]uint32

func Encrypt(key Key, counter uint32, nonce Nonce, data []byte) []byte {
	dataLen := uint32(len(data))
	maxBlock := dataLen / 64
	encrypted := make([]byte, dataLen)
	for i := uint32(0); i <= maxBlock; i++ {
		encryptBlock(key, counter, nonce, data, i, dataLen, encrypted)
	}

	return encrypted
}

func encryptBlock(key Key, counter uint32, nonce Nonce, data []byte, i uint32, dataLen uint32, encrypted []byte) {
	key_stream := blockFunc(key, counter+i, nonce)
	startIndex := i * 64
	finishIndex := startIndex + 64

	if finishIndex > dataLen {
		finishIndex = dataLen
	}

	block := data[startIndex:finishIndex]
	//encryptedBlock := xorArr(key_stream, block)
	xorToArr(key_stream, block, encrypted, int(startIndex))
	//copy(encrypted[startIndex:finishIndex], encryptedBlock)
}

func blockFunc(key Key, counter uint32, nonce Nonce) [64]byte {
	state := createState(key, counter, nonce)
	workingState := state
	workingStatePtr := &workingState

	for i := 0; i < 10; i++ {
		innerBlock(workingStatePtr)
	}

	result := [64]byte{}
	for i := 0; i < len(state); i++ {
		sum := state[i] + workingState[i]
		binary.LittleEndian.PutUint32(result[i*4:], sum)
	}

	return result
}

func createState(key Key, counter uint32, nonce Nonce) [16]uint32 {
	state := [16]uint32{}

	//constants
	state[0] = 0x61707865
	state[1] = 0x3320646e
	state[2] = 0x79622d32
	state[3] = 0x6b206574

	//key
	state[4] = key[0]
	state[5] = key[1]
	state[6] = key[2]
	state[7] = key[3]
	state[8] = key[4]
	state[9] = key[5]
	state[10] = key[6]
	state[11] = key[7]

	//counter, enough for 256GB of data
	state[12] = counter

	//nonce
	state[13] = nonce[0]
	state[14] = nonce[1]
	state[15] = nonce[2]
	return state
}

func xorToArr(a [64]byte, b []byte, result []byte, offset int) {
	for i := 0; i < len(b); i++ {
		result[i+offset] = a[i] ^ b[i]
	}
}

func innerBlock(state *[16]uint32) {
	qRound(state, 0, 4, 8, 12)
	qRound(state, 1, 5, 9, 13)
	qRound(state, 2, 6, 10, 14)
	qRound(state, 3, 7, 11, 15)

	qRound(state, 0, 5, 10, 15)
	qRound(state, 1, 6, 11, 12)
	qRound(state, 2, 7, 8, 13)
	qRound(state, 3, 4, 9, 14)
}

func qRound(state *[16]uint32, a, b, c, d int) {
	//no need for mod max_uint32val, overflow ok
	state[a] += state[b]
	state[d] ^= state[a]
	state[d] = bits.RotateLeft32(state[d], 16)

	state[c] += state[d]
	state[b] ^= state[c]
	state[b] = bits.RotateLeft32(state[b], 12)

	state[a] += state[b]
	state[d] ^= state[a]
	state[d] = bits.RotateLeft32(state[d], 8)

	state[c] += state[d]
	state[b] ^= state[c]
	state[b] = bits.RotateLeft32(state[b], 7)
}
