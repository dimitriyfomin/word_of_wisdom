package protocol

import (
	"context"
	"errors"
	"testing"

	"github.com/dimitriyfomin/word_of_wisdom.git/pkg/commands"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompareBytes(t *testing.T) {
	arr1 := []byte{1, 2, 3}
	arr2 := []byte{1, 2, 3}
	arr3 := []byte{1, 2}
	arr4 := []byte{1, 2, 3, 4}

	assert.True(t, compareBytes(arr1, arr2))

	assert.True(t, compareBytes(nil, nil))

	assert.False(t, compareBytes(arr1, nil))

	assert.False(t, compareBytes(nil, arr2))

	assert.False(t, compareBytes(arr1, arr3))

	assert.False(t, compareBytes(arr1, arr4))
}

func TestParseMessageRequest(t *testing.T) {
	req := ParseMessageRequest([]byte("test"))
	assert.NotNil(t, req)
	assert.Equal(t, "test", req.FullCommandName)
	assert.Nil(t, req.Body)

	req = ParseMessageRequest([]byte("test 000"))
	assert.NotNil(t, req)
	assert.Equal(t, "test", req.FullCommandName)
	assert.NotNil(t, req.Body)

	req = ParseMessageRequest([]byte("test 000 000"))
	assert.NotNil(t, req)
	assert.Equal(t, "test 000", req.FullCommandName)
	assert.NotNil(t, req.Body)
}

func TestNewErrorServerResponse(t *testing.T) {
	req := NewErrorServerResponse(nil)
	assert.NotNil(t, req)
	assert.Equal(t, "ERR Unknown error", string(req))

	req = NewErrorServerResponse(errors.New("New error"))
	assert.NotNil(t, req)
	assert.Equal(t, "ERR New error", string(req))
}

func TestNewOkServerResponse(t *testing.T) {
	req := NewOkServerResponse(nil)
	assert.NotNil(t, req)
	assert.Equal(t, "OK", string(req))

	req = NewOkServerResponse(&commands.MessageResponse{})
	assert.NotNil(t, req)
	assert.Equal(t, "OK", string(req))

	req = NewOkServerResponse(&commands.MessageResponse{
		Body: []byte{},
	})
	assert.NotNil(t, req)
	assert.Equal(t, "OK", string(req))

	req = NewOkServerResponse(&commands.MessageResponse{
		Body: []byte("content"),
	})
	assert.NotNil(t, req)
	assert.Equal(t, "OK content", string(req))
}

func TestHandlePayload(t *testing.T) {
	resp := HandlePayload(nil, nil, nil, nil)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp)
	assert.Equal(t, "ERR Unknown error", string(resp))

	srv := commands.NewServer()
	resp = HandlePayload(srv, nil, []byte("test"), nil)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp)
	assert.Equal(t, "ERR Command test wasn't found", string(resp))

	nothing := func(_ context.Context, _ []byte) ([]byte, error) {
		return nil, nil
	}
	err := srv.Register("", "test", nothing)
	require.NoError(t, err)
	resp = HandlePayload(srv, nil, []byte("test"), nil)
	assert.NotNil(t, resp)
	assert.Equal(t, "OK", string(resp))

	var processedCtx context.Context
	handler := func(ctx context.Context, req []byte) ([]byte, error) {
		processedCtx = ctx
		return req, nil
	}
	err = srv.Register("v1", "echo", handler)
	require.NoError(t, err)
	resp = HandlePayload(srv, nil, []byte("v1.echo test"), nil)
	assert.NotNil(t, resp)
	assert.Equal(t, "OK test", string(resp))
	assert.Nil(t, processedCtx)

	resp = HandlePayload(srv, context.Background(), []byte("v1.echo test"), nil)
	assert.NotNil(t, resp)
	assert.Equal(t, "OK test", string(resp))
	assert.NotNil(t, processedCtx)

	handlerError := func(_ context.Context, _ []byte) ([]byte, error) {
		return nil, errors.New("New error")
	}
	err = srv.Register("v1", "error", handlerError)
	require.NoError(t, err)
	resp = HandlePayload(srv, context.Background(), []byte("v1.error"), nil)
	assert.NotNil(t, resp)
	assert.Equal(t, "ERR New error", string(resp))

	var passedInfo *ExecutionInfo
	interceptor := func(ctx context.Context, req *commands.MessageRequest, info *ExecutionInfo, handler ExecutionHandler) (resp *commands.MessageResponse, err error) {
		passedInfo = info
		return handler(ctx, req)
	}
	resp = HandlePayload(srv, context.Background(), []byte("v1.echo test"), interceptor)
	assert.NotNil(t, resp)
	assert.Equal(t, "OK test", string(resp))
	assert.NotNil(t, passedInfo)
	assert.Equal(t, "v1.echo", passedInfo.CommandName.String())
}
