package node

import (
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/wvalencia19/blocker/proto"
	"github.com/wvalencia19/blocker/types"
)

type BlockStorer interface {
	Put(*proto.Block) error
	Get(string) (*proto.Block, error)
}

type MemoryBlockStore struct {
	lock   sync.RWMutex
	blocks map[string]*proto.Block
}

func NewMemoryBlockStore() *MemoryBlockStore {
	return &MemoryBlockStore{
		blocks: make(map[string]*proto.Block),
	}
}

func (s *MemoryBlockStore) Get(hash string) (*proto.Block, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	block, ok := s.blocks[hash]
	if !ok {
		return nil, fmt.Errorf("block with hash [%s] does not exists", hash)
	}

	return block, nil
}

func (s *MemoryBlockStore) Put(b *proto.Block) error {
	s.lock.RLock()
	defer s.lock.RUnlock()

	hash := hex.EncodeToString(types.HashBlock(b))
	s.blocks[hash] = b
	return nil
}
