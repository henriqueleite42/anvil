package parse

import (
	"errors"
	"strings"

	"github.com/anuntech/hephaestus/cmd/schema"
)

func resolveField(s *schema.Schema, v map[string]any) (map[string]*schema.Field, error) {
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
			entity, ok := s.Entities.Tables[secondPart]
			if !ok {
				return nil, errors.New(refString + ": fail to find entity")
			}

			fields := map[string]*schema.Field{}

			for k, v := range entity.Columns {
				methodField := schema.Field{}

				methodField.Type = v.Type
				methodField.Optional = v.Optional
				methodField.Confidentiality = v.Confidentiality

				fields[k] = &methodField
			}

			return fields, nil
		}

		if firstPart == "Types" {
			if s.Types == nil {
				return nil, errors.New(refString + ": schema has no Types")
			}

			types, ok := (*s.Types)[secondPart]
			if !ok {
				return nil, errors.New(refString + ": fail to find entity")
			}

			return types, nil
		}

		return nil, errors.New("fail to parse something")
	}

	methodField := map[string]*schema.Field{}

	for kk, vv := range v {
		vvMap := vv.(map[string]any)

		var fieldType schema.FieldType
		if val, ok := vvMap["Type"]; ok {
			fieldType = val.(schema.FieldType)
		}
		var dbType *string = nil
		if val, ok := vvMap["DbType"]; ok {
			valString := val.(string)
			dbType = &valString
		}
		var encoded *string = nil
		if val, ok := vvMap["Encoded"]; ok {
			valString := val.(string)
			encoded = &valString
		}
		var optional bool
		if val, ok := vvMap["Optional"]; ok {
			optional = val.(bool)
		}
		var confidentiality schema.FieldConfidentiality = "LOW"
		if val, ok := vvMap["Confidentiality"]; ok {
			confidentiality = val.(schema.FieldConfidentiality)
		}
		var validate []string = nil
		if val, ok := vvMap["Validate"]; ok {
			validate = []string{}
			valSlice := val.([]any)
			for _, v := range valSlice {
				validate = append(validate, v.(string))
			}
		}
		var properties map[string]*schema.Field = nil
		if fieldType == schema.FieldType_Map ||
			fieldType == schema.FieldType_ListMap ||
			fieldType == schema.FieldType_MapStringMap {
			if val, ok := vvMap["Properties"]; ok {
				valMap := val.(map[string]any)
				property, err := resolveField(s, valMap)
				if err != nil {
					return nil, err
				}
				properties = property
			} else {
				return nil, errors.New(kk + " (Map | List[Map] | Map[String]Map) needs Properties")
			}
		}
		var values map[string]string = nil
		if fieldType == schema.FieldType_Enum || fieldType == schema.FieldType_ListEnum {
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

					if s.Enums == nil {
						return nil, errors.New(refString + ": schema has no Enums")
					}
					enumVals, ok := (*s.Enums)[secondPart]
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

		methodField[kk] = &schema.Field{
			Type:            fieldType,
			DbType:          dbType,
			Encoded:         encoded,
			Optional:        optional,
			Confidentiality: confidentiality,
			Validate:        validate,
			Properties:      properties,
			Values:          values,
		}
	}

	return methodField, nil
}

func parseDependency(v map[string]any) (*schema.Dependency, error) {
	depImport := schema.Dependency_Import{}
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

		depImport = schema.Dependency_Import{
			Alias: alias,
			Path:  path,
		}
	}

	var depType string
	if val, ok := v["Type"]; ok {
		valString := val.(string)
		depType = valString
	}

	return &schema.Dependency{
		Type:   depType,
		Import: &depImport,
	}, nil
}
