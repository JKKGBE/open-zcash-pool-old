package proxy

import (
	"encoding/binary"
	"encoding/hex"
)

type extraNonceCounter struct {
	// placeholder []byte
	counter uint32
	// size        int
}

func newExtraNonceCounter(configInstanceId uint32) *extraNonceCounter {
	// placeholder, _ = hex.DecodeString("f000000ff111111f")
	counter := configInstanceId << 27
	// size = 4
	// return &extraNonceCounter{placeholder: placeholder, counter: counter, size: size}
	return &extraNonceCounter{counter: counter}
}

func (n *extraNonceCounter) getNextExtraNonce1() string {
	n.counter += 1
	nonce := make([]byte, 4)
	binary.BigEndian.PutUint32(nonce, n.counter)
	return hex.EncodeToString(nonce)
}

// func (n *extraNonceCounter) getExtraNonce2Size() int {
// 	return len(n.placeholder) - n.size
// }
