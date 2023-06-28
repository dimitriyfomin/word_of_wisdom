package v1

import (
	"github.com/dimitriyfomin/word_of_wisdom.git/pkg/api/wisdom/hub"
	"github.com/dimitriyfomin/word_of_wisdom.git/pkg/commands"
)

const (
	version = "v1"
)

type Api struct {
	Hub hub.Api
}

func (v Api) Register(registrator commands.Registrator) error {
	// The latest version (spike) should be the same as hub

	err := registrator.Register(version, hub.GetChallengeCommand, v.Hub.GetChallenge)
	if err != nil {
		return err
	}

	err = registrator.Register(version, hub.VerifyChallengeCommand, v.Hub.VerifyChallenge)
	if err != nil {
		return err
	}

	err = registrator.Register(version, hub.GetRandomQuoteCommand, v.Hub.GetRandomQuote)
	if err != nil {
		return err
	}

	return nil
}
