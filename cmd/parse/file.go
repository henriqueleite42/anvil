package parse

import "github.com/anuntech/hephaestus/cmd/schema"

func file(data map[string]any) (*schema.Schema, error) {
	schema := schema.Schema{}

	err := domain(&schema, data)
	if err != nil {
		return nil, err
	}

	err = relationships(&schema, data)
	if err != nil {
		return nil, err
	}

	err = types(&schema, data)
	if err != nil {
		return nil, err
	}

	err = enums(&schema, data)
	if err != nil {
		return nil, err
	}

	err = entities(&schema, data)
	if err != nil {
		return nil, err
	}

	err = events(&schema, data)
	if err != nil {
		return nil, err
	}

	err = repository(&schema, data)
	if err != nil {
		return nil, err
	}

	err = usecase(&schema, data)
	if err != nil {
		return nil, err
	}

	err = delivery(&schema, data)
	if err != nil {
		return nil, err
	}

	return &schema, nil
}
