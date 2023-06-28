package hub

import "context"

const (
	challengeContextKey = "challenge"
)

func getContextChallenge(ctx context.Context) []byte {
	if ctx == nil {
		return nil
	}
	challenge := ctx.Value(challengeContextKey)
	if challenge == nil {
		return nil
	}
	return challenge.([]byte)
}

func SetContextChallenge(ctx context.Context, challenge []byte) context.Context {
	if ctx == nil {
		return nil
	}
	return context.WithValue(ctx, challengeContextKey, challenge)
}
