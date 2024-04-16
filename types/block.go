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

func VerifyBlock(b *proto.Block) bool {
	if len(b.PublicKey) != crypto.PubKeyLen {
		return false
	}
	if len(b.Signature) != crypto.SignatureLen {
		return false
	}

	sig := crypto.SignatureFromBytes(b.Signature)
	pubKey := crypto.PublicKeyFromBytes(b.PublicKey)
	hash := HashBlock(b)
	return sig.Verify(pubKey, hash)
}

func SignBlock(pk *crypto.PrivateKey, b *proto.Block) *crypto.Signature {
	hash := HashBlock(b)
	sig := pk.Sign(hash)
	b.PublicKey = pk.Public().Bytes()
	b.Signature = sig.Bytes()
	return sig
}

func HashHeader(header *proto.Header) []byte {
	b, err := pb.Marshal(header)

	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(b)

	return hash[:]
}
