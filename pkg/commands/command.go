package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

type Command struct {
	CommandName CommandName
	handler     CommandHandler
}

type CommandName struct {
	Version string
	Name    string
}

type MessageRequest struct {
	Body []byte
}

type MessageResponse struct {
	Body []byte
}

func checkNamePart(name string) bool {
	for _, r := range []rune{'.', ' '} {
		if strings.ContainsRune(name, r) {
			return false
		}
	}
	return true
}

func NewCommand(version string, name string, handler CommandHandler) (*Command, error) {
	if !checkNamePart(version) {
		return nil, errors.New("Invalid version value")
	}
	if name == "" || !checkNamePart(name) {
		return nil, errors.New("Invalid command name value")
	}
	if handler == nil {
		return nil, errors.New("Invalid handler value")
	}

	return &Command{
		CommandName: CommandName{
			Version: version,
			Name:    name,
		},
		handler: handler,
	}, nil
}

func (cmd CommandName) String() string {
	if cmd.Version == "" {
		return cmd.Name
	}
	return fmt.Sprintf("%s.%s", cmd.Version, cmd.Name)
}

func (cmd Command) Execute(ctx context.Context, req *MessageRequest) (*MessageResponse, error) {
	if cmd.handler == nil {
		return nil, nil
	}

	var body []byte
	if req != nil {
		body = req.Body
	}
	resp, err := cmd.handler(ctx, body)

	return &MessageResponse{
		Body: resp,
	}, err
}
