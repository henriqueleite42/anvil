package templates

type UsecaseStructTemplInput struct {
	Domain      string
	DomainSnake string
}

const UsecaseStructTempl = `package {{ .DomainSnake }}_usecase

type {{ .Domain }}UsecaseImplementation struct {
}
`
