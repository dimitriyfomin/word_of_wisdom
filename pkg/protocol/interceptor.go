package protocol

import (
	"context"

	"github.com/dimitriyfomin/word_of_wisdom.git/pkg/commands"
)

// ExecutionInterceptor provides a hook to intercept the execution of call.
type ExecutionInterceptor func(ctx context.Context, req *commands.MessageRequest, info *ExecutionInfo, handler ExecutionHandler) (resp *commands.MessageResponse, err error)

type ExecutionHandler func(ctx context.Context, req *commands.MessageRequest) (*commands.MessageResponse, error)

type ExecutionInfo struct {
	CommandName commands.CommandName
}
