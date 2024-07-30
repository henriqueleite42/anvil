package parse

import "github.com/anuntech/hephaestus/cmd/types"

func File(data map[string]any) (*types.Schema, error) {
	schema := types.Schema{}

	err := Domain(&schema, data)
	if err != nil {
		return nil, err
	}

	err = Types(&schema, data)
	if err != nil {
		return nil, err
	}

	err = Enums(&schema, data)
	if err != nil {
		return nil, err
	}

	err = Entities(&schema, data)
	if err != nil {
		return nil, err
	}

	err = Events(&schema, data)
	if err != nil {
		return nil, err
	}

	err = Repository(&schema, data)
	if err != nil {
		return nil, err
	}

	err = Usecase(&schema, data)
	if err != nil {
		return nil, err
	}

	err = Delivery(&schema, data)
	if err != nil {
		return nil, err
	}

	return &schema, nil
}
