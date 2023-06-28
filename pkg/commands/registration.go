package commands

import "context"

type CommandHandler func(context.Context, []byte) ([]byte, error)

type Registrator interface {
	Register(version string, commandName string, handler CommandHandler) error
}
