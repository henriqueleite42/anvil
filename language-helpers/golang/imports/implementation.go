package imports

type importsManager struct {
	imports map[string]*Import
}

func NewImportsManager() ImportsManager {
	return &importsManager{
		imports: map[string]*Import{},
	}
}
