package hub

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetChallenge(t *testing.T) {
	api := Api{
		Difficulty: 3,
	}
	challenge, err := api.GetChallenge(nil, nil)

	assert.NoError(t, err)
	assert.NotNil(t, challenge)
}

func TestVerifyChallenge(t *testing.T) {
	challenge, err := hex.DecodeString("5a337512f6c66c1b")
	require.NoError(t, err)
	token, err := hex.DecodeString("3136383735323133333131363131333630373135363237393230303934")
	require.NoError(t, err)
	invalidToken, err := hex.DecodeString("0000000000000000")
	require.NoError(t, err)

	api := Api{
		Difficulty: 3,
	}
	_, err = api.VerifyChallenge(nil, token)
	assert.Error(t, err)

	ctx := context.WithValue(context.Background(), challengeContextKey, challenge)
	resp, err := api.VerifyChallenge(ctx, token)
	assert.NoError(t, err)
	assert.Nil(t, resp)

	_, err = api.VerifyChallenge(ctx, invalidToken)
	assert.Error(t, err)

	ctx = context.WithValue(context.Background(), challengeContextKey, nil)
	resp, err = api.VerifyChallenge(ctx, nil)
	assert.Error(t, err)

	ctx = context.WithValue(context.Background(), challengeContextKey, []byte{})
	resp, err = api.VerifyChallenge(ctx, nil)
	assert.Error(t, err)

	ctx = context.WithValue(context.Background(), challengeContextKey, nil)
	resp, err = api.VerifyChallenge(ctx, []byte{})
	assert.Error(t, err)

	ctx = context.WithValue(context.Background(), challengeContextKey, []byte{})
	resp, err = api.VerifyChallenge(ctx, []byte{})
	assert.Error(t, err)
}
