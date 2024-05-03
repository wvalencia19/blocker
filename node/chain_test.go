package node

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wvalencia19/blocker/crypto"
	"github.com/wvalencia19/blocker/proto"
	"github.com/wvalencia19/blocker/types"
	"github.com/wvalencia19/blocker/util"
)

func RandomBlock(t *testing.T, chain *Chain) *proto.Block {
	privKey := crypto.GeneratePrivatekey()
	b := util.RandomBlock()
	prevBlock, err := chain.GetBlockByHeight(chain.Height())
	require.Nil(t, err)
	b.Header.PrevHash = types.HashBlock(prevBlock)
	types.SignBlock(privKey, b)
	return b
}

func TestNewChain(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTXStore())
	assert.Equal(t, 0, chain.Height())
	_, err := chain.GetBlockByHeight(0)
	assert.Nil(t, err)
}
func TestChainHeight(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTXStore())

	for i := 0; i < 100; i++ {
		b := RandomBlock(t, chain)

		assert.Nil(t, chain.AddBlock(b))
		require.Equal(t, chain.Height(), i+1)
	}

}
func TestAddBlock(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTXStore())
	for i := 0; i < 100; i++ {

		block := RandomBlock(t, chain)
		blockHash := types.HashBlock(block)

		require.Nil(t, chain.AddBlock(block))

		fetchBlock, err := chain.GetBlockByHash(blockHash)
		require.Nil(t, err)
		assert.Equal(t, block, fetchBlock)

		fetchedBlockByHeight, err := chain.GetBlockByHeight(i + 1)
		require.Nil(t, err)
		assert.Equal(t, block, fetchedBlockByHeight)
	}
}

func TestAddBlockWithTx(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTXStore())
	block := RandomBlock(t, chain)
	privKey := crypto.NewPrivateKeyFromSeedStr(godSeed)
	recipient := crypto.GeneratePrivatekey().Public().Address().Bytes()

	prevTx, err := chain.txStore.Get("b78abe0c0dc56af50d070c97bffa92867fda1c26c47455d954533d9f3ce888b6")
	assert.Nil(t, err)

	inputs := []*proto.TxInput{
		{
			PrevTxHash:   types.HashTransaction(prevTx),
			PrevOutIndex: 0,
			PublicKey:    privKey.Public().Bytes(),
		},
	}
	outputs := []*proto.TxOutput{
		{
			Amount:  100,
			Address: recipient,
		},
		{
			Amount:  900,
			Address: privKey.Public().Address().Bytes(),
		},
	}

	tx := &proto.Transaction{
		Version: 1,
		Inputs:  inputs,
		Outputs: outputs,
	}

	sig := types.SignTransaction(privKey, tx)
	tx.Inputs[0].Signature = sig.Bytes()

	block.Transactions = append(block.Transactions, tx)
	require.Nil(t, chain.AddBlock(block))

}

func TestBlockWithTXInsufficientFunds(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTXStore())
	block := RandomBlock(t, chain)
	privKey := crypto.NewPrivateKeyFromSeedStr(godSeed)
	recipient := crypto.GeneratePrivatekey().Public().Address().Bytes()

	prevTx, err := chain.txStore.Get("b78abe0c0dc56af50d070c97bffa92867fda1c26c47455d954533d9f3ce888b6")
	assert.Nil(t, err)

	inputs := []*proto.TxInput{
		{
			PrevTxHash:   types.HashTransaction(prevTx),
			PrevOutIndex: 0,
			PublicKey:    privKey.Public().Bytes(),
		},
	}
	outputs := []*proto.TxOutput{
		{
			Amount:  10001,
			Address: recipient,
		},
	}

	tx := &proto.Transaction{
		Version: 1,
		Inputs:  inputs,
		Outputs: outputs,
	}

	sig := types.SignTransaction(privKey, tx)
	tx.Inputs[0].Signature = sig.Bytes()

	block.Transactions = append(block.Transactions, tx)
	types.SignBlock(privKey, block)
	require.Nil(t, chain.AddBlock(block))
}
