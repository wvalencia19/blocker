package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wvalencia19/blocker/crypto"
	"github.com/wvalencia19/blocker/proto"
	"github.com/wvalencia19/blocker/util"
)

// my balance 100 coins
// want to send 5 coins to "AAA"
// 2 outputs
// 5 to the dude we want to send
// 95 back to our address

func TestNewTransaction(t *testing.T) {
	fromPrivKey := crypto.GeneratePrivatekey()
	fromAddress := fromPrivKey.Public().Address().Bytes()

	toPrivKey := crypto.GeneratePrivatekey()
	toAddress := toPrivKey.Public().Address().Bytes()

	input := &proto.TxInput{
		PrevTxHash:   util.RandomHash(),
		PrevOutIndex: 0,
		PublicKey:    fromPrivKey.Public().Bytes(),
	}

	output1 := &proto.TxOutput{
		Amount:  5,
		Address: toAddress,
	}

	output2 := &proto.TxOutput{
		Amount:  95,
		Address: fromAddress,
	}

	tx := &proto.Transaction{
		Version: 1,
		Inputs:  []*proto.TxInput{input},
		Outputs: []*proto.TxOutput{output1, output2},
	}

	sig := SignTransaction(fromPrivKey, tx)
	input.Signature = sig.Bytes()

	assert.True(t, VerifyTransaction(tx))
}
