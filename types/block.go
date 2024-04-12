package types

import (
	"crypto/sha256"

	"github.com/wvalencia19/blocker/crypto"

	"github.com/wvalencia19/blocker/proto"
	pb "google.golang.org/protobuf/proto"
)

func HashBlock(block *proto.Block) []byte {
	return HashHeader(block.Header)
}

func SignBlock(pk *crypto.PrivateKey, b *proto.Block) *crypto.Signature {
	return pk.Sign(HashBlock(b))
}

func HashHeader(header *proto.Header) []byte {
	b, err := pb.Marshal(header)

	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(b)

	return hash[:]
}
