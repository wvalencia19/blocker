package node

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wvalencia19/blocker/types"
	"github.com/wvalencia19/blocker/util"
)

func TestChainHeight(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore())

	for i := 0; i < 100; i++ {
		b := util.RandomBlock()
		assert.Nil(t, chain.AddBlock(b))
		assert.Equal(t, chain.Height(), i)
	}

}
func TestAddBlock(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore())
	for i := 0; i < 100; i++ {
		var (
			block     = util.RandomBlock()
			blockHash = types.HashBlock(block)
		)
		assert.Nil(t, chain.AddBlock(block))

		fetchBlock, err := chain.GetBlockByHash(blockHash)
		assert.Nil(t, err)
		assert.Equal(t, block, fetchBlock)

		fetchedBlockByHeight, err := chain.GetBlockByHeight(i)
		assert.Nil(t, err)
		assert.Equal(t, block, fetchedBlockByHeight)
	}
}
