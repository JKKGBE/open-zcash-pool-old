package proxy

import "github.com/cbergoon/merkletree"

func getRoot(transactions transact, generatedTxRaw string) []byte {
	if len(transactions) == 0 {
		return generatedTxRaw
	}
		
	var hashes [len(transactions) + 1]string
	hashes = 

	t, err := merkletree.NewTree(hashes)
	if err != nil {
		fmt.Prinln("Error creating Merkle Tree:", err)
	}
	return t.MerkleRoot
}
