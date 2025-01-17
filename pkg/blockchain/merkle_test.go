package blockchain

import (
	"testing"

	block2 "github.com/cybriq/p9/pkg/block"
)

// TestMerkle tests the BuildMerkleTreeStore API.
func TestMerkle(t *testing.T) {
	block := block2.NewBlock(&Block100000)
	merkles := BuildMerkleTreeStore(block.Transactions(), false)
	calculatedMerkleRoot := merkles.GetRoot()
	wantMerkle := &Block100000.Header.MerkleRoot
	if !wantMerkle.IsEqual(calculatedMerkleRoot) {
		t.Errorf(
			"BuildMerkleTreeStore: merkle root mismatch - got %v, want %v",
			calculatedMerkleRoot, wantMerkle,
		)
	}
}
