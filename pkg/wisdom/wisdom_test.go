package wisdom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRandomQuote(t *testing.T) {
	quote := GetRandomQuote()

	assert.NotEqual(t, quote, "")
}
