package hub

import (
	"github.com/dimitriyfomin/word_of_wisdom.git/pkg/commands"
)

const (
	version                = "hub"
	GetChallengeCommand    = "CHG"
	VerifyChallengeCommand = "CHGT"
	GetRandomQuoteCommand  = "QTR"
)

type Api struct {
	Difficulty int
}

func (h Api) Register(registrator commands.Registrator) error {
	err := registrator.Register(version, GetChallengeCommand, h.GetChallenge)
	if err != nil {
		return err
	}

	err = registrator.Register(version, VerifyChallengeCommand, h.VerifyChallenge)
	if err != nil {
		return err
	}

	err = registrator.Register(version, GetRandomQuoteCommand, h.GetRandomQuote)
	if err != nil {
		return err
	}

	return nil
}
