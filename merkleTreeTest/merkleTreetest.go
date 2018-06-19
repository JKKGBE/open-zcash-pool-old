package main

import (
	"fmt"
	"strconv"
)

func main() {
	fmt.Println("elo")
	fmt.Println(strconv.FormatInt(32, 16))
}

// func main() {
// 	b1, _ := hex.DecodeString("999d2c8bb6bda0bf784d9ebeb631d711dbbbfe1bc006ea13d6ad0d6a2649a971")
// 	b2, _ := hex.DecodeString("3f92594d5a3d7b4df29d7dd7c46a0dac39a96e751ba0fc9bab5435ea5e22a19d")
// 	b3, _ := hex.DecodeString("a5633f03855f541d8e60a6340fc491d49709dc821f3acb571956a856637adcb6")
// 	// b4, _ := hex.DecodeString("28d97c850eaf917a4c76c02474b05b70a197eaefb468d21c22ed110afe8ec9e0")

// 	var b1a [32]byte
// 	var b2a [32]byte
// 	var b3a [32]byte

// 	copy(b1a[:], b1)
// 	copy(b2a[:], b2)
// 	copy(b3a[:], b3)

// 	concat1 := b1a[:]
// 	concat1 = append(concat1, b2a[:]...)

// 	concat2 := b3a[:]
// 	concat2 = append(concat2, b3a[:]...)

// 	// fmt.Println(b1a)
// 	// fmt.Println(b2a)
// 	// fmt.Println(concat1)

// 	shaConcat1 := sha256.Sum256(concat1)
// 	fmt.Println(shaConcat1)
// 	// hexShaConcat1 := hex.EncodeToString(shaConcat1[:])
// 	// fmt.Println(hexShaConcat1)

// 	shaConcat1 = sha256.Sum256(shaConcat1[:])
// 	fmt.Println(shaConcat1)
// 	// hexShaConcat1 = hex.EncodeToString(shaConcat1[:])
// 	// fmt.Println(hexShaConcat1)

// 	shaConcat2 := sha256.Sum256(concat2)
// 	shaConcat2 = sha256.Sum256(shaConcat2[:])
// 	// hexShaConcat2 := hex.EncodeToString(shaConcat2[:])

// 	concat3 := shaConcat1[:]
// 	concat3 = append(concat3, shaConcat2[:]...)

// 	shaConcat3 := sha256.Sum256(concat3)
// 	shaConcat3 = sha256.Sum256(shaConcat3[:])
// 	hexShaConcat3 := hex.EncodeToString(shaConcat3[:])
// 	fmt.Println(hexShaConcat3)
// }

// func main() {
// 	b1, _ := hex.DecodeString("999d2c8bb6bda0bf784d9ebeb631d711dbbbfe1bc006ea13d6ad0d6a2649a971")
// 	b2, _ := hex.DecodeString("3f92594d5a3d7b4df29d7dd7c46a0dac39a96e751ba0fc9bab5435ea5e22a19d")
// 	b3, _ := hex.DecodeString("a5633f03855f541d8e60a6340fc491d49709dc821f3acb571956a856637adcb6")
// 	// b4, _ := hex.DecodeString("28d97c850eaf917a4c76c02474b05b70a197eaefb468d21c22ed110afe8ec9e0")
// 	// 	fmt.Println(b1)
// 	var b1a [32]byte
// 	var b2a [32]byte
// 	var b3a [32]byte
// 	// var b4a [32]byte
// 	copy(b1a[:], b1)
// 	copy(b2a[:], b2)
// 	copy(b3a[:], b3)
// 	// copy(b4a[:], b4)

// 	bottomRow := merkleTree.Row{b1a, b2a, b3a}

// 	mt := merkleTree.NewMerkleTree(bottomRow)
// 	mr := mt.MerkleRoot()
// 	root := hex.EncodeToString(mr[:])
// 	fmt.Println(root)

// 	// TODO test, expected hex value for 4 tx is
// 	// "fa64c69618475f3c7bf8ff64bd3a1be37d80690822b81c5b7301a9e7800f498e"
// 	//
// 	// and for 3 tx is
// 	// "0a8ce36d84bccedf18dcc900ded1131ebbe233d67c1ed9ae7ad31a1030c18d4f"
// }
