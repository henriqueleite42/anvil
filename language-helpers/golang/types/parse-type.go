package types_parser

import (
	"fmt"
	"slices"
	"strings"

	"github.com/henriqueleite42/anvil/language-helpers/golang/imports"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

// originalRootNode is used so we can pass to the child types their parent module
// and see if it need to be imported or not
func (self *typeParser) ParseType(t *schemas.Type) (*Type, error) {
	result := &Type{
		Optional:  t.Optional,
		AnvilType: t,
	}

	// ----------------------
	//
	// Basic types
	//
	// ----------------------

	if t.Type == schemas.TypeType_String {
		result.GolangType = "string"
	}
	if t.Type == schemas.TypeType_Int {
		result.GolangType = "int"
	}
	if t.Type == schemas.TypeType_Int8 {
		result.GolangType = "int8"
	}
	if t.Type == schemas.TypeType_Int16 {
		result.GolangType = "int16"
	}
	if t.Type == schemas.TypeType_Int32 {
		result.GolangType = "int32"
	}
	if t.Type == schemas.TypeType_Int64 {
		result.GolangType = "int64"
	}
	if t.Type == schemas.TypeType_Uint {
		result.GolangType = "uint"
	}
	if t.Type == schemas.TypeType_Uint8 {
		result.GolangType = "uint8"
	}
	if t.Type == schemas.TypeType_Uint16 {
		result.GolangType = "uint16"
	}
	if t.Type == schemas.TypeType_Uint32 {
		result.GolangType = "uint32"
	}
	if t.Type == schemas.TypeType_Uint64 {
		result.GolangType = "uint64"
	}
	if t.Type == schemas.TypeType_Float {
		result.GolangType = "float"
	}
	if t.Type == schemas.TypeType_Float32 {
		result.GolangType = "float32"
	}
	if t.Type == schemas.TypeType_Float64 {
		result.GolangType = "float64"
	}
	if t.Type == schemas.TypeType_Bool {
		result.GolangType = "bool"
	}
	if t.Type == schemas.TypeType_Timestamp {
		result.GolangType = "Time"
		result.ModuleImport = imports.NewImport("time", nil)
		result.imports = imports.NewImportsManager()
		result.imports.AddImport("time", nil)
	}

	// ----------------------
	//
	// Complex types
	//
	// ----------------------

	if t.Type == schemas.TypeType_Enum {
		schemaEnum := self.schema.Enums.Enums[*t.EnumHash]
		enum, err := self.ParseEnum(schemaEnum)
		if err != nil {
			return nil, err
		}

		result.GolangType = enum.GolangName
		result.ModuleImport = self.getEnumsImport(schemaEnum)
		result.imports = imports.NewImportsManager()
		result.imports.MergeImport(result.ModuleImport)
	}
	if t.Type == schemas.TypeType_List {
		childType, ok := self.schema.Types.Types[t.ChildTypes[0].TypeHash]
		if !ok {
			return nil, fmt.Errorf("[types: type enum] type \"%s\" not found", t.ChildTypes[0].TypeHash)
		}

		resolvedChildType, err := self.ParseType(childType)
		if err != nil {
			return nil, err
		}

		var golangType string
		if resolvedChildType.AnvilType.Type == schemas.TypeType_Map ||
			resolvedChildType.Optional {
			golangType = "[]*" + resolvedChildType.GolangType
		} else {
			golangType = "[]" + resolvedChildType.GolangType
		}

		result.GolangType = golangType
		result.ModuleImport = resolvedChildType.ModuleImport
		result.imports = resolvedChildType.imports
	}

	// ----------------------
	//
	// Maps (child types)
	//
	// ----------------------

	if t.Type == schemas.TypeType_Map {
		if existentType, ok := self.typesToAvoidDuplication[t.Ref]; ok {
			return existentType, nil
		}

		props := make([]*MapProp, len(t.ChildTypes))

		result.imports = imports.NewImportsManager()

		for k, v := range t.ChildTypes {
			if v.PropName == nil {
				return nil, fmt.Errorf("ChildType \"%s.%d\" must have a PropName", t.Name, k)
			}

			childType, ok := self.schema.Types.Types[v.TypeHash]
			if !ok {
				return nil, fmt.Errorf("[types: type map] type \"%s\" not found", v.TypeHash)
			}

			propType, err := self.ParseType(childType)
			if err != nil {
				return nil, err
			}

			if propType.ModuleImport != nil {
				// Maps should import all of their properties modules
				result.imports.MergeImport(propType.ModuleImport)
			}

			prop := &MapProp{
				Name: *v.PropName,
				Type: propType,
			}

			// Very complicated and counterintuitive, we need to do this.
			// Learn more in:
			// https://pkg.go.dev/github.com/go-playground/validator/v10#hdr-Omit_Empty
			// https://pkg.go.dev/github.com/go-playground/validator/v10#hdr-Required
			if childType.Optional {
				if !slices.Contains(childType.Validate, "omitempty") {
					// Needs to be the first thing on the slice
					childType.Validate = append([]string{"omitempty"}, childType.Validate...)
				}
			} else {
				if !slices.Contains(childType.Validate, "required") &&
					// Only applies it to pointer types, because when you use the "required"
					// tag on "go-playground/validator", it only checks if the value is != than it's zero
					// value, what makes booleans fail if false and numbers fail if 0,
					// what is not the expected intuitive behavior of "required" (it should
					// be checking if the value was received, and not it's contents)
					!IsTypeBool(childType) {
					// Needs to be the first thing on the slice
					childType.Validate = append([]string{"required"}, childType.Validate...)
				}
			}

			if len(childType.Validate) > 0 || len(childType.Transform) > 0 || childType.DbName != nil {
				if prop.Tags == nil {
					prop.Tags = []string{}
				}

				if len(childType.Validate) > 0 {
					prop.Tags = append(prop.Tags, fmt.Sprintf("validate:\"%s\"", strings.Join(childType.Validate, ",")))
				}

				if len(childType.Transform) > 0 {
					prop.Tags = append(prop.Tags, fmt.Sprintf("mod:\"%s\"", strings.Join(childType.Transform, ",")))
				}

				if childType.DbName != nil {
					prop.Tags = append(prop.Tags, fmt.Sprintf("db:\"%s\"", *childType.DbName))
				}
			}

			props[k] = prop
		}

		result.GolangType = t.Name
		result.MapProps = props

		if t.RootNode == "Types" {
			result.ModuleImport = self.getTypesImport(t)
			result.imports.AddImport(result.ModuleImport.Path, &result.ModuleImport.Alias)
			self.types = append(self.types, result)
		} else if t.RootNode == "Events" {
			result.ModuleImport = self.getEventsImport(t)
			result.imports.AddImport(result.ModuleImport.Path, &result.ModuleImport.Alias)
			self.events = append(self.events, result)
		} else if t.RootNode == "Entities" {
			result.ModuleImport = self.getEntitiesImport(t)
			result.imports.AddImport(result.ModuleImport.Path, &result.ModuleImport.Alias)
			self.entities = append(self.entities, result)
		} else if t.RootNode == "Repository" {
			result.ModuleImport = self.getRepositoryImport(t)
			result.imports.AddImport(result.ModuleImport.Path, &result.ModuleImport.Alias)
			self.repository = append(self.repository, result)
		} else if t.RootNode == "Usecase" {
			result.ModuleImport = self.getUsecaseImport(t)
			result.imports.AddImport(result.ModuleImport.Path, &result.ModuleImport.Alias)
			self.usecase = append(self.usecase, result)
		} else {
			return nil, fmt.Errorf("unable to get package for \"%s\"", t.Ref)
		}
	}

	self.typesToAvoidDuplication[t.Ref] = result

	return result, nil
}
