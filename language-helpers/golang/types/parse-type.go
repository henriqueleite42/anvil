package types_parser

import (
	"fmt"
	"slices"
	"strings"

	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func (self *typeParser) getNodePkg(node string) (string, error) {
	if node == "Enums" {
		return self.enumsPkg, nil
	} else if node == "Types" {
		return self.typesPkg, nil
	} else if node == "Events" {
		return self.eventsPkg, nil
	} else if node == "Entities" {
		return self.entitiesPkg, nil
	} else if node == "Repository" {
		return self.repositoryPkg, nil
	} else if node == "Usecase" {
		return self.usecasePkg, nil
	} else {
		return "", fmt.Errorf("unable to get package  node \"%s\"", node)
	}
}

// originalRootNode is used so we can pass to the child types their parent module
// and see if it need to be imported or not
func (self *typeParser) parseType(t *schemas.Type, originalRootNode string) (*Type, error) {
	if originalRootNode == "" {
		originalRootNode = t.RootNode
	}

	result := &Type{
		Optional:  t.Optional,
		AnvilType: t.Type,
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
		result.GolangType = "time.Time"

		if t.RootNode == "Types" {
			self.AddTypesImport("time")
		} else if t.RootNode == "Events" {
			self.AddEventsImport("time")
		} else if t.RootNode == "Entities" {
			self.AddEntitiesImport("time")
		} else if t.RootNode == "Repository" {
			self.AddRepositoryImport("time")
		} else if t.RootNode == "Usecase" {
			self.AddUsecaseImport("time")
		} else {
			return nil, fmt.Errorf("unable to get package for \"%s\"", t.Ref)
		}
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

		result.GolangPkg = &enum.GolangPkg
		result.GolangType = enum.GolangName

		if schemaEnum.RootNode != originalRootNode {
			if originalRootNode == "Types" {
				self.AddTypesImport(self.enumsPkg)
			} else if originalRootNode == "Events" {
				self.AddEventsImport(self.enumsPkg)
			} else if originalRootNode == "Entities" {
				self.AddEntitiesImport(self.enumsPkg)
			} else if originalRootNode == "Repository" {
				self.AddRepositoryImport(self.enumsPkg)
			} else if originalRootNode == "Usecase" {
				self.AddUsecaseImport(self.enumsPkg)
			} else {
				return nil, fmt.Errorf("unable to get package for \"%s\"", t.Ref)
			}
		}
	}
	if t.Type == schemas.TypeType_List {
		childType, ok := self.schema.Types.Types[t.ChildTypes[0].TypeHash]
		if !ok {
			return nil, fmt.Errorf("type \"%s\" not found", t.ChildTypes[0].TypeHash)
		}

		resolvedChildType, err := self.parseType(childType, t.RootNode)
		if err != nil {
			return nil, err
		}

		var golangType string
		if resolvedChildType.AnvilType == schemas.TypeType_Map ||
			resolvedChildType.Optional {
			golangType = "[]*" + resolvedChildType.GolangType
		} else {
			golangType = "[]" + resolvedChildType.GolangType
		}

		if childType.RootNode != originalRootNode {
			pkgToImport, err := self.getNodePkg(childType.RootNode)
			if err != nil {
				return nil, err
			}

			if originalRootNode == "Types" {
				self.AddTypesImport(pkgToImport)
			} else if originalRootNode == "Events" {
				self.AddEventsImport(pkgToImport)
			} else if originalRootNode == "Entities" {
				self.AddEntitiesImport(pkgToImport)
			} else if originalRootNode == "Repository" {
				self.AddRepositoryImport(pkgToImport)
			} else if originalRootNode == "Usecase" {
				self.AddUsecaseImport(pkgToImport)
			} else {
				return nil, fmt.Errorf("unable to get package for \"%s\"", t.Ref)
			}
		}

		result.GolangPkg = resolvedChildType.GolangPkg
		result.GolangType = golangType
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

		for k, v := range t.ChildTypes {
			if v.PropName == nil {
				return nil, fmt.Errorf("ChildType \"%s.%d\" must have a PropName", t.Name, k)
			}

			childType, ok := self.schema.Types.Types[v.TypeHash]
			if !ok {
				return nil, fmt.Errorf("type \"%s\" not found", v.TypeHash)
			}

			if childType.RootNode != originalRootNode {
				pkgToImport, err := self.getNodePkg(childType.RootNode)
				if err != nil {
					return nil, err
				}

				if originalRootNode == "Types" {
					self.AddTypesImport(pkgToImport)
				} else if originalRootNode == "Events" {
					self.AddEventsImport(pkgToImport)
				} else if originalRootNode == "Entities" {
					self.AddEntitiesImport(pkgToImport)
				} else if originalRootNode == "Repository" {
					self.AddRepositoryImport(pkgToImport)
				} else if originalRootNode == "Usecase" {
					self.AddUsecaseImport(pkgToImport)
				} else {
					return nil, fmt.Errorf("unable to get package for \"%s\"", t.Ref)
				}
			}

			propType, err := self.parseType(childType, childType.RootNode)
			if err != nil {
				return nil, err
			}

			prop := &MapProp{
				Name: *v.PropName,
				Type: propType,
			}

			if !childType.Optional && !slices.Contains(childType.Validate, "required") {
				childType.Validate = append(childType.Validate, "required")
			}

			if len(childType.Validate) > 0 {
				if prop.Tags == nil {
					prop.Tags = []string{}
				}

				prop.Tags = append(prop.Tags, fmt.Sprintf("validate:\"%s\"", strings.Join(childType.Validate, ",")))
			}

			if childType.DbName != nil {
				prop.Tags = append(prop.Tags, fmt.Sprintf("db:\"%s\"", *childType.DbName))
			}

			props[k] = prop
		}

		if len(props) > 0 {
			biggestName := 0
			biggestType := 0
			for _, v := range props {
				newLenName := len(v.Name)
				if newLenName > biggestName {
					biggestName = newLenName
				}

				newLenType := len(v.Type.GolangType)
				if newLenType > biggestType {
					biggestType = newLenType
				}
			}

			for _, v := range props {
				targetLenName := biggestName - len(v.Name)
				v.Spacing1 = strings.Repeat(" ", targetLenName)

				targetLenType := biggestType - len(v.Type.GolangType)
				v.Spacing2 = strings.Repeat(" ", targetLenType)
			}
		}

		result.GolangType = t.Name
		result.MapProps = props
	}

	if t.RootNode == "Types" {
		result.GolangPkg = &self.typesPkg
		self.types = append(self.types, result)
	} else if t.RootNode == "Events" {
		result.GolangPkg = &self.eventsPkg
		self.events = append(self.events, result)
	} else if t.RootNode == "Entities" {
		result.GolangPkg = &self.entitiesPkg
		self.entities = append(self.entities, result)
	} else if t.RootNode == "Repository" {
		result.GolangPkg = &self.repositoryPkg
		self.repository = append(self.repository, result)
	} else if t.RootNode == "Usecase" {
		result.GolangPkg = &self.usecasePkg
		self.usecase = append(self.usecase, result)
	} else {
		return nil, fmt.Errorf("unable to get package for \"%s\"", t.Ref)
	}

	self.typesToAvoidDuplication[t.Ref] = result

	return result, nil
}

func (self *typeParser) ParseType(t *schemas.Type) (*Type, error) {
	return self.parseType(t, "")
}
