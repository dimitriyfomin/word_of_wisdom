package commands

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	srv := NewServer()
	assert.NotNil(t, srv)
	assert.NotNil(t, srv.commands)
}

func TestServerGetCommand(t *testing.T) {
	srv := NewServer()
	srv.commands["v1.test"] = &Command{
		CommandName: CommandName{
			Version: "v1",
			Name:    "test",
		},
	}
	cmd, err := srv.GetCommand("v1.test")
	assert.NoError(t, err)
	assert.NotNil(t, cmd)

	cmd, err = srv.GetCommand("invalid")
	assert.Error(t, err)
	assert.Nil(t, cmd)

	srv.commands = nil
	cmd, err = srv.GetCommand("v1.test")
	assert.Error(t, err)
	assert.Nil(t, cmd)
}

func TestServerRegister(t *testing.T) {
	handler := func(_ context.Context, payload []byte) ([]byte, error) {
		return nil, nil
	}
	srv := NewServer()
	err := srv.Register("v1", "test", handler)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(srv.commands))

	err = srv.Register("v1", "test", handler)
	assert.Error(t, err)
	assert.Equal(t, 1, len(srv.commands))

	err = srv.Register("v1", "test2", nil)
	assert.Error(t, err)
	assert.Equal(t, 1, len(srv.commands))

	srv.commands = nil
	err = srv.Register("v1", "cmd", handler)
	assert.Error(t, err)
	assert.Nil(t, srv.commands)
}
