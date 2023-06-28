package hub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRandomQuote(t *testing.T) {
	api := Api{}
	quote, err := api.GetRandomQuote(nil, nil)

	assert.NoError(t, err)
	assert.NotNil(t, quote)
}
