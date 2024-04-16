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

func TestSignVerifyBlock(t *testing.T) {
	block := util.RandomBlock()
	privKey := crypto.GeneratePrivatekey()
	pubKey := privKey.Public()

	sig := SignBlock(privKey, block)
	assert.Equal(t, 64, len(sig.Bytes()))
	assert.True(t, sig.Verify(pubKey, HashBlock(block)))

	assert.Equal(t, block.PublicKey, pubKey.Bytes())
	assert.Equal(t, block.Signature, sig.Bytes())

	assert.True(t, VerifyBlock(block))

	invalidPrivKey := crypto.GeneratePrivatekey()
	block.PublicKey = invalidPrivKey.Public().Bytes()
	assert.False(t, VerifyBlock(block))
}
