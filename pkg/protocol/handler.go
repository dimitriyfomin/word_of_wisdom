package protocol

import (
	"context"
	"fmt"

	"github.com/dimitriyfomin/word_of_wisdom.git/pkg/commands"
)

const (
	ErrorCommand = "ERR"
	OkCommand    = "OK"
)

var bodyDelimiterBytes = []byte(" ")

type CommandStorage interface {
	GetCommand(string) (*commands.Command, error)
}

func HandlePayload(srv CommandStorage, ctx context.Context, payload []byte, interceptor ExecutionInterceptor) []byte {
	if srv == nil {
		return NewErrorServerResponse(nil)
	}

	req := ParseMessageRequest(payload)
	cmd, err := srv.GetCommand(req.FullCommandName)
	if err != nil || cmd == nil {
		return NewErrorServerResponse(err)
	}

	var resp *commands.MessageResponse
	if interceptor != nil {
		info := &ExecutionInfo{
			CommandName: cmd.CommandName,
		}
		resp, err = interceptor(ctx, req, info, cmd.Execute)
	} else {
		resp, err = cmd.Execute(ctx, req)
	}

	if err != nil {
		return NewErrorServerResponse(err)
	}
	return NewOkServerResponse(resp)
}

// Message format: <fullCommand> <payload:optional>
func ParseMessageRequest(msg []byte) *commands.MessageRequest {
	var fullCommandName string
	var body []byte
	for ind := range msg {
		if compareBytes(msg[ind:ind+len(bodyDelimiterBytes)], bodyDelimiterBytes) {
			fullCommandName = string(msg[:ind])
			body = msg[ind+len(bodyDelimiterBytes):]
		}
	}
	if fullCommandName == "" {
		fullCommandName = string(msg)
	}
	return &commands.MessageRequest{
		FullCommandName: fullCommandName,
		Body:            body,
	}
}

// Message format: ERR <message>
func NewErrorServerResponse(err error) []byte {
	if err == nil {
		return []byte(fmt.Sprintf("%s Unknown error", ErrorCommand))
	}
	return []byte(fmt.Sprintf("%s %v", ErrorCommand, err))
}

// Message format: OK <payload:optional>
func NewOkServerResponse(resp *commands.MessageResponse) []byte {
	if resp == nil || resp.Body == nil || len(resp.Body) == 0 {
		return []byte(OkCommand)
	}

	return append([]byte(OkCommand+" "), resp.Body...)
}

func compareBytes(arr1 []byte, arr2 []byte) bool {
	if arr1 == nil && arr2 == nil {
		return true
	}
	if arr1 == nil || arr2 == nil {
		return false
	}
	if len(arr1) != len(arr2) {
		return false
	}
	for ind, v := range arr1 {
		if v != arr2[ind] {
			return false
		}
	}
	return true
}
