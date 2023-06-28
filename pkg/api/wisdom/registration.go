package wisdom

import (
	"github.com/dimitriyfomin/word_of_wisdom.git/pkg/api/wisdom/hub"
	v1 "github.com/dimitriyfomin/word_of_wisdom.git/pkg/api/wisdom/v1"
	"github.com/dimitriyfomin/word_of_wisdom.git/pkg/commands"
)

func Register(registrator commands.Registrator, difficulty int) error {
	apiHub := hub.Api{
		Difficulty: difficulty,
	}
	err := apiHub.Register(registrator)
	if err != nil {
		return err
	}

	apiV1 := v1.Api{Hub: apiHub}
	err = apiV1.Register(registrator)
	if err != nil {
		return err
	}

	return nil
}
