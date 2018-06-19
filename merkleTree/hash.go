package merkleTree

import (
	"crypto/sha256"
	"fmt"
)

// [32]byte is a type that provides a cute way of expressing this trivial
// fixed size array type, and is useful because the type is often used in
// slices, and it prevents a forest of square brackets when that is done.
// type [32]byte [32]byte

// Hash is a trivial wrapper around one of Go's native hashing functions, that
// serves only to make calls to it look simpler and declutter from them the
// specifics of which hash variant is being used.
// Bitcoin does mandate the use of this particular hashing algorithm, but
// requires in most cases that it be applied twice. To have done so in our
// example code would have added needless complexity.
func DoubleHash(input []byte) [32]byte {
	firstHashing := sha256.Sum256(input)
	return sha256.Sum256(firstHashing[:])
}

// JoinAndHash is a function of fundamental importance to this example code
// because it is the hashing function to derive parents in MerkleTrees from the
// hash values of their children. It follows the Bitcoin specification in that
// it concatentates the two given hashes (as byte streams) and re-hashes the
// result. (But only once).
func JoinAndHash(left [32]byte, right [32]byte) [32]byte {
	combined := left[:]
	combined = append(combined, right[:]...)
	return DoubleHash(combined)
}

// The Hex method is just syntax sugar to avoid having to write things like
// fmt.Sprintf("0x", ... all over the place.
func Hex(h [32]byte) string {
	return fmt.Sprintf("%0x", h)
}
