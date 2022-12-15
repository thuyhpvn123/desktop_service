package bls

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestNewKeyPair(t *testing.T) {
	keyPair := NewKeyPair(common.FromHex("4b1d80795e7404e9bc20fb5d0e783400dbf29ef006ad9518c712a8c60d33c2dd"))
	assert.NotNil(t, keyPair)
}
