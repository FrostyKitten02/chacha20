package chacha20

import (
	"ChaCha20/internal"
	"encoding/binary"
	"errors"
)

func GetKeyFromStr(key string) (*internal.Key, error) {
	if len(key) != 32 {
		return nil, errors.New("invalid key length")
	}

	return bytesToKey([]byte(key)), nil
}

func bytesToKey(val []byte) *internal.Key {
	key := &internal.Key{}
	for i := 0; i < 8; i++ {
		key[i] = binary.LittleEndian.Uint32(val[i*4 : i*4+4])
	}

	return key
}

func GetNonceFromString(nonce string) (*internal.Nonce, error) {
	if len(nonce) != 12 {
		return nil, errors.New("invalid nonce length")
	}

	return bytesToNonce([]byte(nonce)), nil
}

func bytesToNonce(val []byte) *internal.Nonce {
	nonce := &internal.Nonce{}
	for i := 0; i < 3; i++ {
		nonce[i] = binary.LittleEndian.Uint32(val[i*4 : i*4+4])
	}

	return nonce
}
