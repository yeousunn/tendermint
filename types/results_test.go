package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestABCIResults(t *testing.T) {
	a := ABCIResult{Code: 0, Data: nil}
	b := ABCIResult{Code: 0, Data: []byte{}}
	c := ABCIResult{Code: 0, Data: []byte("one")}
	d := ABCIResult{Code: 14, Data: nil}
	e := ABCIResult{Code: 14, Data: []byte("foo")}
	f := ABCIResult{Code: 14, Data: []byte("bar")}

	// Nil and []byte{} should produce the same hash.
	require.Equal(t, a.Hash(), a.Hash())
	require.Equal(t, b.Hash(), b.Hash())
	require.Equal(t, a.Hash(), b.Hash())

	// a and b should be the same, don't go in results.
	results := ABCIResults{a, c, d, e, f}

	// Make sure each result hashes properly.
	var last []byte
	for i, res := range results {
		h := res.Hash()
		assert.NotEqual(t, last, h, "%d", i)
		last = h
	}

	// Make sure that we can get a root hash from results and verify proofs.
	root := results.Hash()
	assert.NotEmpty(t, root)

	for i, res := range results {
		proof := results.ProveResult(i)
		valid := proof.Verify(root, res.Hash())
		assert.NoError(t, valid, "%d", i)
	}
}

func TestABCIResultsBytes(t *testing.T) {
	results := NewResults([]*abci.ResponseDeliverTx{
		{Code: 0, Data: []byte{}},
		{Code: 0, Data: []byte("one")},
		{Code: 14, Data: nil},
		{Code: 14, Data: []byte("foo")},
		{Code: 14, Data: []byte("bar")},
	})
	assert.NotNil(t, results.Bytes())
}
