package parse

import (
	"errors"
	"strings"

	"github.com/anuntech/hephaestus/cmd/types"
)

func resolveField(schema *types.Schema, v map[string]any) (map[string]*types.Field, error) {
	if ref, ok := v["$ref"]; ok {
		refString, ok := ref.(string)
		if !ok {
			return nil, errors.New("fail to parse something")
		}

		parts := strings.Split(refString, ".")
		if len(parts) != 2 {
			return nil, errors.New(refString + ": $ref must have 2 parts divided in 2 by a dot")
		}

		firstPart := parts[0]
		secondPart := parts[1]

		if firstPart == "Entities" {
			entity, ok := schema.Entities.Tables[secondPart]
			if !ok {
				return nil, errors.New(refString + ": fail to find entity")
			}

			fields := map[string]*types.Field{}

			for k, v := range entity.Columns {
				methodField := types.Field{}

				methodField.Type = v.Type
				methodField.Optional = v.Optional
				methodField.Confidentiality = v.Confidentiality

				fields[k] = &methodField
			}

			return fields, nil
		}

		if firstPart == "Types" {
			if schema.Types == nil {
				return nil, errors.New(refString + ": schema has no Types")
			}

			types, ok := (*schema.Types)[secondPart]
			if !ok {
				return nil, errors.New(refString + ": fail to find entity")
			}

			return types, nil
		}

		return nil, errors.New("fail to parse something")
	}

	methodField := map[string]*types.Field{}

	for kk, vv := range v {
		vvMap := vv.(map[string]any)

		var fieldType types.FieldType
		if val, ok := vvMap["Type"]; ok {
			fieldType = val.(types.FieldType)
		}
		var dbType *string = nil
		if val, ok := vvMap["DbType"]; ok {
			valString := val.(string)
			dbType = &valString
		}
		var optional bool
		if val, ok := vvMap["Optional"]; ok {
			optional = val.(bool)
		}
		var confidentiality types.FieldConfidentiality = "LOW"
		if val, ok := vvMap["Confidentiality"]; ok {
			confidentiality = val.(types.FieldConfidentiality)
		}
		var validate []string = nil
		if val, ok := vvMap["Validate"]; ok {
			validate = []string{}
			valSlice := val.([]any)
			for _, v := range valSlice {
				validate = append(validate, v.(string))
			}
		}
		var properties map[string]*types.Field = nil
		if fieldType == types.FieldType_Map {
			if val, ok := vvMap["Properties"]; ok {
				valMap := val.(map[string]any)
				property, err := resolveField(schema, valMap)
				if err != nil {
					return nil, err
				}
				properties = property
			} else {
				return nil, errors.New(kk + " (Map) needs Properties")
			}
		}
		var items map[string]*types.Field = nil
		if fieldType == types.FieldType_List {
			if val, ok := vvMap["Items"]; ok {
				valMap := val.(map[string]any)
				item, err := resolveField(schema, valMap)
				if err != nil {
					return nil, err
				}
				items = item
			} else {
				return nil, errors.New(kk + " (List) needs Properties")
			}
		}
		var values map[string]string = nil
		if fieldType == types.FieldType_Enum || fieldType == types.FieldType_ListEnum {
			if val, ok := vvMap["Values"]; ok {
				valMap := val.(map[string]any)

				ref, ok := valMap["$ref"]
				if ok {
					refString := ref.(string)

					parts := strings.Split(refString, ".")
					if len(parts) != 2 {
						return nil, errors.New(refString + ": $ref must have 2 parts divided in 2 by a dot")
					}

					secondPart := parts[1]

					if schema.Enums == nil {
						return nil, errors.New(refString + ": schema has no Enums")
					}
					enumVals, ok := (*schema.Enums)[secondPart]
					if !ok {
						return nil, errors.New(refString + ": fail to find entity")
					}

					values = enumVals
				} else {
					enumVals := map[string]string{}
					for kkk, vvv := range valMap {
						enumVals[kkk] = vvv.(string)
					}
					values = enumVals
				}
			} else {
				return nil, errors.New(kk + " (Enum | List[Enum]) needs Values")
			}
		}

		methodField[kk] = &types.Field{
			Type:            fieldType,
			DbType:          dbType,
			Optional:        optional,
			Confidentiality: confidentiality,
			Validate:        validate,
			Properties:      properties,
			Items:           items,
			Values:          values,
		}
	}

	return methodField, nil
}

func parseDependency(v map[string]any) (*types.Dependency, error) {
	depImport := types.Dependency_Import{}
	depImportAny, ok := v["Import"]
	if ok {
		depImportMap := depImportAny.(map[string]any)

		var alias *string = nil
		if val, ok := depImportMap["Alias"]; ok {
			valString := val.(string)
			alias = &valString
		}

		var path string
		if val, ok := depImportMap["Path"]; ok {
			valString := val.(string)
			path = valString
		}

		depImport = types.Dependency_Import{
			Alias: alias,
			Path:  path,
		}
	}

	var depType string
	if val, ok := v["Type"]; ok {
		valString := val.(string)
		depType = valString
	}

	return &types.Dependency{
		Type:   depType,
		Import: &depImport,
	}, nil
}
