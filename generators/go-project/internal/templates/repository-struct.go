package templates

type RepositoryStructTemplInput struct {
	Domain      string
	DomainSnake string
}

const RepositoryStructTempl = `package {{ .DomainSnake }}_repository

type {{ .Domain }}RepositoryImplementation struct {
}
`
