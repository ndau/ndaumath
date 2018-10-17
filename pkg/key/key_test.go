package key

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMaster1(t *testing.T) {
	k, err := NewMaster([]byte("abcdefghijklmnopqrstuvwxyz123456"))
	assert.Nil(t, err)
	fmt.Println(k)
}

func TestChildren(t *testing.T) {
	k, err := NewMaster([]byte("abcdefghijklmnopqrstuvwxyz123456"))
	assert.Nil(t, err)
	ch0, err := k.Child(0)
	assert.Nil(t, err)
	_, err = k.Child(23456)
	assert.Nil(t, err)
	ch0a, err := k.Child(0)
	assert.Nil(t, err)
	assert.Equal(t, ch0, ch0a)
}

func checkKeys(t *testing.T, pvt, pub *ExtendedKey) {
	pvtk, err := pvt.ECPrivKey()
	assert.Nil(t, err)
	pubk, err := pub.ECPubKey()
	assert.Nil(t, err)
	h := doubleHashB([]byte("This is a test message to hash"))
	sig, err := pvtk.Sign(h)
	assert.Nil(t, err)
	assert.True(t, sig.Verify(h, pubk))
}

func TestGenPublicBasic(t *testing.T) {
	pvt, err := NewMaster([]byte("abcdefghijklmnopqrstuvwxyz123456"))
	assert.Nil(t, err)
	pub, err := pvt.Public()
	assert.Nil(t, err)
	checkKeys(t, pvt, pub)
}

func TestGenPublicChild(t *testing.T) {
	// Check that we can generate the first child from a private key, then
	// derive a public key from it, and then sign/verify using those keys
	pvtmaster, err := NewMaster([]byte("abcdefghijklmnopqrstuvwxyz123456"))
	assert.Nil(t, err)
	pvt, err := pvtmaster.Child(0)
	assert.Nil(t, err)
	pub, err := pvt.Public()
	assert.Nil(t, err)
	checkKeys(t, pvt, pub)
}

func TestGenPublicChild2(t *testing.T) {
	// Check that we can generate a high-numbered child from a private key, then
	// derive a public key from it, and then sign/verify using those keys
	pvtmaster, err := NewMaster([]byte("abcdefghijklmnopqrstuvwxyz123456"))
	assert.Nil(t, err)
	pvt, err := pvtmaster.Child(400)
	assert.Nil(t, err)
	pub, err := pvt.Public()
	assert.Nil(t, err)
	checkKeys(t, pvt, pub)
}

func TestGenPublicChild3(t *testing.T) {
	// Check that we can generate the 4th child from a private key, then
	// derive a public key from the parent key, and then generate the 4th child from that,
	// and then sign/verify using those keys
	pvtmaster, err := NewMaster([]byte("abcdefghijklmnopqrstuvwxyz123456"))
	assert.Nil(t, err)
	pvt, err := pvtmaster.Child(4)
	assert.Nil(t, err)
	pubmaster, err := pvtmaster.Public()
	assert.Nil(t, err)
	pub, err := pubmaster.Child(4)
	assert.Nil(t, err)
	checkKeys(t, pvt, pub)
}

func TestPubPrv(t *testing.T) {
	pvtmaster, err := NewMaster([]byte("abcdefghijklmnopqrstuvwxyz123456"))
	assert.Nil(t, err)
	pvt, err := pvtmaster.Child(1)
	assert.Nil(t, err)
	checkKeys(t, pvt, pvt)
}
