package hub

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetContextChallenge(t *testing.T) {
	challenge, err := hex.DecodeString("5a337512f6c66c1b")
	require.NoError(t, err)

	assert.Nil(t, getContextChallenge(nil))

	ctx := context.Background()
	assert.Nil(t, getContextChallenge(ctx))

	ctx = context.WithValue(ctx, challengeContextKey, challenge)
	res := getContextChallenge(ctx)
	assert.NotNil(t, res)
	assert.Equal(t, "5a337512f6c66c1b", hex.EncodeToString(getContextChallenge(ctx)))
}

func TestSetContextChallenge(t *testing.T) {
	challenge, err := hex.DecodeString("5a337512f6c66c1b")
	require.NoError(t, err)

	ctx := SetContextChallenge(nil, challenge)
	assert.Nil(t, ctx)

	ctx = context.Background()
	ctx = SetContextChallenge(ctx, challenge)
	assert.NotNil(t, ctx)
	assert.Equal(t, "5a337512f6c66c1b", hex.EncodeToString(ctx.Value(challengeContextKey).([]byte)))
}
