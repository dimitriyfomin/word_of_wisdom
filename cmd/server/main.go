package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"time"

	"github.com/dimitriyfomin/word_of_wisdom.git/pkg/api/wisdom"
	"github.com/dimitriyfomin/word_of_wisdom.git/pkg/api/wisdom/hub"
	"github.com/dimitriyfomin/word_of_wisdom.git/pkg/commands"
	"github.com/dimitriyfomin/word_of_wisdom.git/pkg/protocol"
)

const difficulty = 3 // We set the difficulty of the PoW

type ConnectionValues struct {
	ddosSecureChecked bool
	challenge         []byte
}

func main() {
	rand.Seed(time.Now().UnixNano())

	commandsServer := commands.NewServer()
	err := wisdom.Register(commandsServer, difficulty)
	if err != nil {
		panic(err)
	}

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	ctx := context.Background()

	fmt.Println("Server is running!")
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn, commandsServer, ctx)
	}
}

func handleConnection(conn net.Conn, srv *commands.Server, ctx context.Context) {
	defer conn.Close()

	connValues := ConnectionValues{}
	interceptor := func(ctx context.Context, req *commands.MessageRequest, info *protocol.ExecutionInfo, handler protocol.ExecutionHandler) (*commands.MessageResponse, error) {
		return interceptBySecurity(ctx, req, info, handler, &connValues)
	}

	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err == io.EOF {
			return
		} else if err != nil {
			fmt.Println("Error reading request:", err)
			return
		}

		responseBytes := protocol.HandlePayload(srv, ctx, buffer[:n], interceptor)

		_, err = conn.Write(responseBytes)
		if err != nil {
			fmt.Println("Error writing response:", err)
			return
		}
	}
}

func interceptBySecurity(ctx context.Context, req *commands.MessageRequest, info *protocol.ExecutionInfo, handler protocol.ExecutionHandler, connValues *ConnectionValues) (*commands.MessageResponse, error) {
	// Require DDoS check
	if connValues != nil && !connValues.ddosSecureChecked {
		ignoreCheck := false
		if info != nil {
			for _, cmdName := range []string{hub.VerifyChallengeCommand, hub.GetChallengeCommand} {
				if cmdName == info.CommandName.Name {
					ignoreCheck = true
				}
			}
		}
		if !ignoreCheck {
			return nil, errors.New(fmt.Sprintf("Can't call command %s until anti-DDoS check", info.CommandName.String()))
		}
	}

	// Challenge enrichment for context
	if info != nil && connValues != nil && info.CommandName.Name == hub.VerifyChallengeCommand {
		ctx = hub.SetContextChallenge(ctx, connValues.challenge)
	}

	// Default processing in handler
	var resp *commands.MessageResponse
	var err error
	if handler != nil {
		resp, err = handler(ctx, req)
	}

	// Store challenge after generation
	if info != nil && connValues != nil && info.CommandName.Name == hub.GetChallengeCommand && resp != nil && err == nil {
		connValues.challenge = resp.Body
	}
	// DDoS Secure check has been passed
	if info != nil && connValues != nil && info.CommandName.Name == hub.VerifyChallengeCommand && err == nil {
		connValues.ddosSecureChecked = true
	}

	return resp, err
}
