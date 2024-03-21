package types

import (
	"testing"

	"github.com/wvalencia19/blocker/crypto"

	"github.com/stretchr/testify/assert"
	"github.com/wvalencia19/blocker/util" // Import the util package
)

func TestHashBlock(t *testing.T) {
	block := util.RandomBlock()
	hash := HashBlock(block)
	assert.Equal(t, 32, len(hash))
}

func TestSignBlock(t *testing.T) {
	block := util.RandomBlock()

	privKey := crypto.GeneratePrivatekey()
	pubKey := privKey.Public()

	sig := SignBlock(privKey, block)
	assert.Equal(t, 64, len(sig.Bytes()))
	assert.True(t, sig.Verify(pubKey, HashBlock(block)))
}
