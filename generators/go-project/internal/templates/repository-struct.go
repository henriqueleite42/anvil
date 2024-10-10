package templates

const RepositoryStructTempl = `package {{ .DomainSnake }}_repository

import (
	"errors"

	"github.com/rs/zerolog"
)

type {{ .DomainCamel }}RepositoryImplementation struct {
	logger *zerolog.Logger
}

type New{{ .Domain }}RepositoryInput struct {
	Logger *zerolog.Logger
}

func New{{ .Domain }}Repository(i *New{{ .Domain }}RepositoryInput) ({{ .Domain }}Repository, error) {
	if i == nil {
		return nil, errors.New("New{{ .Domain }}Repository: input must not be nil")
	}

	return &{{ .DomainCamel }}RepositoryImplementation{
		logger: i.Logger,
	}, nil
}
`
