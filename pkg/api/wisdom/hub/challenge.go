package hub

import (
	"context"
	"errors"

	"github.com/dimitriyfomin/word_of_wisdom.git/pkg/security"
)

func (h *Api) GetChallenge(ctx context.Context, _ []byte) ([]byte, error) {
	challenge, err := security.GenerateChallenge(h.Difficulty)
	if err != nil {
		return nil, err
	}
	return challenge, err
}

func (h *Api) VerifyChallenge(ctx context.Context, payload []byte) ([]byte, error) {
	challenge := getContextChallenge(ctx)
	if challenge == nil {
		return nil, errors.New("No challenge generated")
	}

	if security.VerifyPoW(challenge, payload, h.Difficulty) {
		return nil, nil
	} else {
		return nil, errors.New("Verification failed")
	}
}
