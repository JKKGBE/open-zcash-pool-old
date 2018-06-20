package transactions

// import (
// 	"bytes"
// 	"encoding/binary"
// 	"io"

// 	"github.com/btcsuite/btcd/chaincfg/chainhash"
// 	"github.com/btcsuite/btcd/wire"
// )

// type Tx struct {
// 	IsOverwinter       bool
// 	Version            uint32
// 	VersionGroupID     uint32
// 	Inputs             []Input
// 	Outputs            []Output
// 	LockTime           uint32
// 	ExpiryHeight       uint32
// 	JoinSplits         []JoinSplit
// 	JoinSplitPubKey    [32]byte
// 	JoinSplitSignature [64]byte
// }

// type Input struct {
// 	PreviousOutPoint wire.OutPoint
// 	SignatureScript  []byte
// 	Sequence         uint32
// }

// type Output struct {
// 	Value        int64
// 	ScriptPubKey []byte
// }

// type JoinSplit struct {
// 	VPubOld      uint64
// 	VPubNew      uint64
// 	Anchor       [32]byte
// 	Nullifiers   [2][32]byte
// 	Commitments  [2][32]byte
// 	EphemeralKey [32]byte
// 	RandomSeed   [32]byte
// 	Macs         [2][32]byte
// 	Proof        [296]byte
// 	Ciphertexts  [2][601]byte
// }

// type countingWriter struct {
// 	io.Writer
// 	N int64
// }

// func CreateRawTransaction(inputs []Input, outputs []Output) *Tx {
// 	tx := &Tx{
// 		IsOverwinter:   false,
// 		Version:        1,
// 		VersionGroupID: 0,
// 		Inputs:         inputs,
// 		Outputs:        outputs,
// 		LockTime:       0,
// 	}

// 	if tx.IsOverwinter {
// 		tx.Version = 3
// 		tx.VersionGroupID = 0x03C48270
// 	}

// 	return tx
// }

// func (t *Tx) GetHeader() uint32 {
// 	if t.IsOverwinter {
// 		return t.Version | OverwinterFlagMask
// 	}
// 	return t.Version
// }

// func (t *Tx) TxHash() chainhash.Hash {
// 	b, _ := t.MarshalBinary()
// 	return chainhash.DoubleHashH(b)
// }

// func (t *Tx) MarshalBinary() ([]byte, error) {
// 	buf := &bytes.Buffer{}
// 	if _, err := t.WriteTo(buf); err != nil {
// 		return nil, err
// 	}
// 	return buf.Bytes(), nil
// }

// func (t *Tx) WriteTo(w io.Writer) (n int64, err error) {
// 	counter := &countingWriter{Writer: w}
// 	for _, segment := range []func(io.Writer) error{
// 		writeField(t.GetHeader()),
// 		writeIf(t.IsOverwinter, writeField(t.VersionGroupID)),
// 		t.writeInputs,
// 		t.writeOutputs,
// 		writeField(t.LockTime),
// 		writeIf(t.IsOverwinter, writeField(t.ExpiryHeight)),
// 		writeIf(t.Version >= 2, t.writeJoinSplits),
// 		writeIf(t.Version >= 2 && len(t.JoinSplits) > 0, writeBytes(t.JoinSplitPubKey[:])),
// 		writeIf(t.Version >= 2 && len(t.JoinSplits) > 0, writeBytes(t.JoinSplitSignature[:])),
// 	} {
// 		if err := segment(counter); err != nil {
// 			return counter.N, err
// 		}
// 	}
// 	return counter.N, nil
// }

// func writeField(v interface{}) func(w io.Writer) error {
// 	return func(w io.Writer) error {
// 		return binary.Write(w, binary.LittleEndian, v)
// 	}
// }

// func writeIf(pred bool, f func(w io.Writer) error) func(w io.Writer) error {
// 	if pred {
// 		return f
// 	}
// 	return func(w io.Writer) error { return nil }
// }

// func writeBytes(v []byte) func(w io.Writer) error {
// 	return func(w io.Writer) error {
// 		_, err := w.Write(v)
// 		return err
// 	}
// }
