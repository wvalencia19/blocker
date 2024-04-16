package node

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/wvalencia19/blocker/crypto"
	"github.com/wvalencia19/blocker/proto"
	"github.com/wvalencia19/blocker/types"
)

type HeaderList struct {
	headers []*proto.Header
}

func NewHeaderList() *HeaderList {
	return &HeaderList{
		headers: []*proto.Header{},
	}
}
func (list *HeaderList) Add(h *proto.Header) {
	list.headers = append(list.headers, h)
}

func (list *HeaderList) Get(index int) *proto.Header {
	if index > list.Height() {
		panic("index too high")
	}
	return list.headers[index]
}

func (list *HeaderList) Height() int {
	return list.Len() - 1
}

// [A, B, C, D, E] len = 5, height = 4
func (list *HeaderList) Len() int {
	return len(list.headers)
}

type Chain struct {
	blockStore BlockStorer
	headers    *HeaderList
}

func NewChain(bs BlockStorer) *Chain {
	chain := &Chain{
		blockStore: bs,
		headers:    NewHeaderList(),
	}
	chain.addBlock(CreateGenesisBlock())

	return chain
}

func (c *Chain) Height() int {
	return c.headers.Height()
}

func (c *Chain) AddBlock(b *proto.Block) error {
	if err := c.ValidateBlock(b); err != nil {
		return err
	}
	// add the header to the list of headers
	return c.addBlock(b)
}

func (c *Chain) addBlock(b *proto.Block) error {
	c.headers.Add(b.Header)
	// validation
	return c.blockStore.Put(b)
}

func (c *Chain) GetBlockByHeight(height int) (*proto.Block, error) {
	if c.Height() < height {
		return nil, fmt.Errorf("given height [%d] too high - height [%d]", height, c.Height())
	}

	header := c.headers.Get(height)
	hash := types.HashHeader(header)
	return c.GetBlockByHash(hash)
}

func (c *Chain) GetBlockByHash(hash []byte) (*proto.Block, error) {
	hashHex := hex.EncodeToString(hash)
	return c.blockStore.Get(hashHex)
}

func (c *Chain) ValidateBlock(b *proto.Block) error {
	/// validate the signature of the block
	if !types.VerifyBlock(b) {
		return fmt.Errorf("invalid block signature")
	}
	// validate if the prevHash is the actually hash of the current block
	currentBlock, err := c.GetBlockByHeight((c.Height()))
	if err != nil {
		return err
	}
	hash := types.HashBlock(currentBlock)
	if !bytes.Equal(hash, b.Header.PrevHash) {
		return fmt.Errorf("invalid previous block hash")
	}
	return nil
}

func CreateGenesisBlock() *proto.Block {
	privKey := crypto.GeneratePrivatekey()
	block := &proto.Block{
		Header: &proto.Header{
			Version: 1,
		},
	}
	types.SignBlock(privKey, block)
	return block
}
