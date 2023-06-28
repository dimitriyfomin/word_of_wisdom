package security

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const difficulty = 3
const attempts = 100_000_000

func TestVerifyPoW(t *testing.T) {
	challenge, err := hex.DecodeString("5a337512f6c66c1b")
	require.NoError(t, err)
	token, err := hex.DecodeString("3136383735323133333131363131333630373135363237393230303934")
	require.NoError(t, err)

	result := VerifyPoW(challenge, token, difficulty)
	assert.True(t, result)
}

func TestGenerateChallenge(t *testing.T) {
	challenge, err := GenerateChallenge(difficulty)

	assert.NoError(t, err)
	assert.True(t, len(challenge) > 0)
}

func TestGenerateTokenByChallenge(t *testing.T) {
	challenge, err := hex.DecodeString("5a337512f6c66c1b")
	assert.NoError(t, err)

	token := GenerateTokenByChallenge(challenge, difficulty, attempts)

	assert.NotNil(t, token)
	assert.True(t, len(token) > 0, token)
}
