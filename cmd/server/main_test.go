package main

import (
	"context"
	"errors"
	"testing"

	"github.com/dimitriyfomin/word_of_wisdom.git/pkg/api/wisdom/hub"
	"github.com/dimitriyfomin/word_of_wisdom.git/pkg/commands"
	"github.com/dimitriyfomin/word_of_wisdom.git/pkg/protocol"
	"github.com/stretchr/testify/assert"
)

func TestInterceptBySecurity(t *testing.T) {
	resp, err := interceptBySecurity(nil, nil, nil, nil, nil)
	assert.NoError(t, err)
	assert.Nil(t, resp)

	ctx := context.Background()
	execInfo := &protocol.ExecutionInfo{
		CommandName: commands.CommandName{
			Version: "",
			Name:    "test",
		},
	}
	emptyReq := &commands.MessageRequest{Body: nil}
	values := &ConnectionValues{}
	_, err = interceptBySecurity(ctx, emptyReq, execInfo, nil, values)
	assert.Error(t, err, "No access rights until ddos check")

	challengeReq := &commands.MessageRequest{Body: []byte("challenge")}
	execInfo = &protocol.ExecutionInfo{
		CommandName: commands.CommandName{
			Version: "",
			Name:    hub.GetChallengeCommand,
		},
	}
	echoHandler := func(_ context.Context, req *commands.MessageRequest) (*commands.MessageResponse, error) {
		if req == nil {
			return nil, nil
		}
		return &commands.MessageResponse{
			Body: req.Body,
		}, nil
	}
	resp, err = interceptBySecurity(ctx, challengeReq, execInfo, echoHandler, values)
	assert.NoError(t, err, "GetChallenge should not be covered by ddos check")
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Body)
	assert.Equal(t, "challenge", string(resp.Body))
	assert.NotNil(t, values.challenge)
	assert.Equal(t, "challenge", string(values.challenge))
	assert.False(t, values.ddosSecureChecked)

	execInfo = &protocol.ExecutionInfo{
		CommandName: commands.CommandName{
			Version: "",
			Name:    hub.VerifyChallengeCommand,
		},
	}
	errorHandler := func(_ context.Context, _ *commands.MessageRequest) (*commands.MessageResponse, error) {
		return nil, errors.New("New error")
	}
	resp, err = interceptBySecurity(ctx, emptyReq, execInfo, errorHandler, values)
	assert.Error(t, err, "VerifyChallenge should not be passed by error check")
	assert.Nil(t, resp)
	assert.NotNil(t, values.challenge)
	assert.False(t, values.ddosSecureChecked)

	execInfo = &protocol.ExecutionInfo{
		CommandName: commands.CommandName{
			Version: "",
			Name:    hub.VerifyChallengeCommand,
		},
	}
	resp, err = interceptBySecurity(ctx, emptyReq, execInfo, echoHandler, values)
	assert.NoError(t, err, "VerifyChallenge should not be covered by ddos check")
	assert.NotNil(t, resp)
	assert.Nil(t, resp.Body)
	assert.NotNil(t, values.challenge)
	assert.True(t, values.ddosSecureChecked)

	quoteReq := &commands.MessageRequest{Body: []byte("quote")}
	execInfo = &protocol.ExecutionInfo{
		CommandName: commands.CommandName{
			Version: "",
			Name:    "test",
		},
	}
	resp, err = interceptBySecurity(ctx, quoteReq, execInfo, echoHandler, values)
	assert.NoError(t, err, "No ddos for any call check when security passed")
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Body)
	assert.Equal(t, "quote", string(resp.Body))
}
