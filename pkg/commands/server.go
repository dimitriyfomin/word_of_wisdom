package commands

import (
	"errors"
	"fmt"
)

type Server struct {
	commands map[string]*Command
}

func NewServer() *Server {
	return &Server{
		commands: make(map[string]*Command),
	}
}

func (srv Server) GetCommand(fullCommandName string) (*Command, error) {
	if srv.commands == nil {
		return nil, errors.New("Commands weren't initialized")
	}

	if cmd, ok := srv.commands[fullCommandName]; ok {
		return cmd, nil
	}
	return nil, errors.New(fmt.Sprintf("Command %s wasn't found", fullCommandName))
}

func (srv *Server) Register(version string, commandName string, handler CommandHandler) error {
	cmd, err := NewCommand(version, commandName, handler)
	if err != nil {
		return err
	}
	if cmd == nil {
		return errors.New("Can't create command")
	}
	fullCommand := cmd.CommandName.String()

	if srv.commands == nil {
		return errors.New("Commands weren't initialized")
	}

	if _, ok := srv.commands[fullCommand]; ok {
		return errors.New(fmt.Sprintf("Can't register command %s twice", fullCommand))
	}
	srv.commands[fullCommand] = &Command{
		CommandName: CommandName{
			Version: version,
			Name:    commandName,
		},
		handler: handler,
	}

	return nil
}
