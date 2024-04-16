package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivatekey(t *testing.T) {
	privKey := GeneratePrivatekey()
	assert.Equal(t, len(privKey.Bytes()), PrivKeyLen)

	pubKey := privKey.Public()
	assert.Equal(t, len(pubKey.Bytes()), PubKeyLen)

}

func TestPrivateKeySign(t *testing.T) {
	privKey := GeneratePrivatekey()
	pubKey := privKey.Public()
	msg := []byte("foo bar baz")

	sig := privKey.Sign(msg)
	assert.True(t, sig.Verify(pubKey, msg))
	assert.False(t, sig.Verify(pubKey, []byte("foo")))

	invalidPrivKey := GeneratePrivatekey()
	invalidPubKey := invalidPrivKey.Public()
	assert.False(t, sig.Verify(invalidPubKey, msg))
}

func TestPublicKeyToAddress(t *testing.T) {
	privKey := GeneratePrivatekey()
	pubKey := privKey.Public()
	address := pubKey.Address()

	assert.Equal(t, AddressLen, len(address.Bytes()))
}

func TestNewPrivateKeyFromString(t *testing.T) {
	seed := "e8482210a5ae3c338733e7b124849c8e7fd350e01bdd017e0eb83bd16815b39e"
	addressStr := "5fe9a8a54115a3a404e1405bd3d9a162961dcac4"
	privKey := NewPrivateKeyFromString(seed)
	assert.Equal(t, PrivKeyLen, len(privKey.Bytes()))

	address := privKey.Public().Address()
	assert.Equal(t, addressStr, address.String())
}
