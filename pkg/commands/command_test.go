package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCommand(t *testing.T) {
	handler := func(_ context.Context, _ []byte) ([]byte, error) {
		return nil, nil
	}
	cmd, err := NewCommand("v1", "test", handler)
	assert.NoError(t, err)
	assert.Equal(t, "v1", cmd.CommandName.Version)
	assert.Equal(t, "test", cmd.CommandName.Name)
	assert.NotNil(t, cmd.handler)

	_, err = NewCommand("", "test", handler)
	assert.NoError(t, err)

	_, err = NewCommand("v1.v2", "test", handler)
	assert.Error(t, err)

	_, err = NewCommand("v1 v2", "test", handler)
	assert.Error(t, err)

	_, err = NewCommand("v1", "", handler)
	assert.Error(t, err)

	_, err = NewCommand("v1", "test.another", handler)
	assert.Error(t, err)

	_, err = NewCommand("v1", "test another", handler)
	assert.Error(t, err)

	_, err = NewCommand("v1", "test", nil)
	assert.Error(t, err)
}

func TestCommandString(t *testing.T) {
	handler := func(_ context.Context, _ []byte) ([]byte, error) {
		return nil, nil
	}
	cmd, err := NewCommand("v1", "test", handler)
	assert.NoError(t, err)
	assert.Equal(t, "v1.test", cmd.CommandName.String())

	cmd, err = NewCommand("", "test", handler)
	assert.NoError(t, err)
	assert.Equal(t, "test", cmd.CommandName.String())
}

func TestCommandExecute(t *testing.T) {
	processed := false
	var req []byte
	handler := func(_ context.Context, payload []byte) ([]byte, error) {
		processed = true
		req = payload
		return nil, nil
	}
	cmd, err := NewCommand("v1", "test", handler)
	_, err = cmd.Execute(nil, nil)
	assert.NoError(t, err)
	assert.True(t, processed)
	assert.Nil(t, req)

	processed = false
	cmd, err = NewCommand("v1", "test", handler)
	_, err = cmd.Execute(nil, &MessageRequest{Body: []byte{0, 0}})
	assert.NoError(t, err)
	assert.True(t, processed)
	assert.NotNil(t, req)
}
