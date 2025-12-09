package blockchain

import (
	"crypto/sha256"
)

// BuildMerkleRoot builds a Merkle root from a slice of transaction hashes.
// If the number of leaves is odd, the last leaf is duplicated.
func BuildMerkleRoot(leaves [][]byte) []byte {
	if len(leaves) == 0 {
		return []byte{}
	}

	// If only one leaf, return its hash (or hash of itself to be consistent)
	if len(leaves) == 1 {
		h := sha256.Sum256(leaves[0])
		return h[:]
	}

	nodes := make([][]byte, len(leaves))
	copy(nodes, leaves)

	for len(nodes) > 1 {
		var nextLevel [][]byte
		for i := 0; i < len(nodes); i += 2 {
			if i+1 == len(nodes) {
				// duplicate last
				data := append(nodes[i], nodes[i]...)
				h := sha256.Sum256(data)
				nextLevel = append(nextLevel, h[:])
			} else {
				data := append(nodes[i], nodes[i+1]...)
				h := sha256.Sum256(data)
				nextLevel = append(nextLevel, h[:])
			}
		}
		nodes = nextLevel
	}

	return nodes[0]
}
