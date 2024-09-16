package parser

import (
	"fmt"
	"strings"

	"github.com/henriqueleite42/anvil/cli/schemas"
	"github.com/henriqueleite42/anvil/generators/go-project/internal/templates"
)

type ResolveMapPropInput struct {
	Type              *schemas.Type
	Kind              Kind
	PrefixForChildren string
	Tags              []string
}

func (self *Parser) ResolveMapProp(i *ResolveMapPropInput) (*templates.TemplTypeProp, error) {
	result := &templates.TemplTypeProp{
		Name: i.Type.Name,
	}

	if i.Type.Type == schemas.TypeType_String {
		result.Type = "string"
	}
	if i.Type.Type == schemas.TypeType_Int {
		result.Type = "int32"
	}
	if i.Type.Type == schemas.TypeType_Float {
		result.Type = "float32"
	}
	if i.Type.Type == schemas.TypeType_Bool {
		result.Type = "bool"
	}
	if i.Type.Type == schemas.TypeType_Timestamp {
		if i.Kind == Kind_Entity || i.Kind == Kind_Event {
			self.ImportsModels["time"] = true
		} else if i.Kind == Kind_Repository {
			self.ImportsRepository["time"] = true
		} else if i.Kind == Kind_Usecase {
			self.ImportsUsecase["time"] = true
		}
		result.Type = "time.Time"
	}
	if i.Type.Type == schemas.TypeType_Enum {
		if i.Type.EnumHash == nil {
			return nil, fmt.Errorf("enum for type \"%s\" not found", i.Type.Name)
		}
		schemaEnum, ok := self.Schema.Enums.Enums[*i.Type.EnumHash]
		if !ok {
			return nil, fmt.Errorf("enum \"%s\" of type \"%s\" not found", *i.Type.EnumHash, i.Type.Name)
		}

		e, err := self.resolveEnum(schemaEnum)
		if err != nil {
			return nil, err
		}

		if i.Kind == Kind_Repository {
			self.ImportsRepository[self.ModelsPath] = true
		} else if i.Kind == Kind_Usecase {
			self.ImportsUsecase[self.ModelsPath] = true
		}

		if i.Kind == Kind_Repository || i.Kind == Kind_Usecase {
			result.Type = fmt.Sprintf("%s.%s", self.ModelsPkgName, e.Name)
		} else {
			// In models file
			result.Type = e.Name
		}
	}
	if i.Type.Type == schemas.TypeType_List {
		if i.Type.ChildTypesHashes == nil {
			return nil, fmt.Errorf("ChildTypesHashes for \"%s\" not found", i.Type.Name)
		}
		if len(i.Type.ChildTypesHashes) != 1 {
			return nil, fmt.Errorf("ChildTypesHashes for \"%s\" must have exactly one item", i.Type.Name)
		}

		childType, ok := self.Schema.Types.Types[i.Type.ChildTypesHashes[0]]
		if !ok {
			return nil, fmt.Errorf("type \"%s\" not found", i.Type.ChildTypesHashes[0])
		}

		if childType.Type == schemas.TypeType_Map {
			resolvedChildType, err := self.ResolveMap(i.Kind, childType, i.PrefixForChildren)
			if err != nil {
				return nil, err
			}

			result.Type = "[]*" + resolvedChildType.Name
		} else {
			resolvedChildType, err := self.ResolveMapProp(&ResolveMapPropInput{
				Kind:              i.Kind,
				Type:              childType,
				PrefixForChildren: i.PrefixForChildren,
			})
			if err != nil {
				return nil, err
			}

			result.Type = "[]" + resolvedChildType.Type
		}
	}
	if i.Type.Type == schemas.TypeType_Map {
		resolvedChildType, err := self.ResolveMap(i.Kind, i.Type, i.PrefixForChildren)
		if err != nil {
			return nil, err
		}

		result.Type = "*" + resolvedChildType.Name
	}

	if i.Type.Optional && i.Type.Type != schemas.TypeType_Map && i.Type.Type != schemas.TypeType_List {
		result.Type = "*" + result.Type
	}

	if i.Kind == Kind_Usecase && len(i.Type.Validate) > 0 {
		i.Tags = append(i.Tags, fmt.Sprintf("validate:\"%s\"", strings.Join(i.Type.Validate, ",")))
	}

	if len(i.Tags) > 0 {
		result.Tags = fmt.Sprintf("`%s`", strings.Join(i.Tags, " "))
	}

	return result, nil
}

func (self *Parser) ResolveMap(kind Kind, t *schemas.Type, prefix string) (*templates.TemplType, error) {
	if t == nil {
		return nil, fmt.Errorf("type must be specified, received nil")
	}

	name := prefix + t.Name

	if kind == Kind_Repository {
		if existent, ok := self.TypesRepositoryToAvoidDuplication[name]; ok {
			return existent, nil
		}
	} else if kind == Kind_Usecase {
		if existent, ok := self.TypesUsecaseToAvoidDuplication[name]; ok {
			return existent, nil
		}
	}

	if self.Schema == nil {
		return nil, fmt.Errorf("missing schema")
	}
	if self.Schema.Types == nil || self.Schema.Types.Types == nil {
		return nil, fmt.Errorf("missing schema types")
	}
	if t.Type != schemas.TypeType_Map {
		return nil, fmt.Errorf("type \"%s\" must be a map, received \"%s\"", t.Name, t.Type)
	}
	if t.ChildTypesHashes == nil {
		return nil, fmt.Errorf("type \"%s\" missing required field \"ChildTypesHashes\"", t.Name)
	}
	lenChildTypesHashes := len(t.ChildTypesHashes)
	if lenChildTypesHashes == 0 {
		return nil, fmt.Errorf("type \"%s\" has to have at least 1 \"ChildTypesHashes\"", t.Name)
	}

	result := &templates.TemplType{
		Name:         name,
		OriginalType: t.Type,
		Props:        make([]*templates.TemplTypeProp, lenChildTypesHashes, lenChildTypesHashes),
	}

	for k, v := range t.ChildTypesHashes {
		childTypeRef, ok := self.Schema.Types.Types[v]
		if !ok {
			return nil, fmt.Errorf("child type \"%s\" of type \"%s\" not found", v, t.Name)
		}

		resolvedProp, err := self.ResolveMapProp(&ResolveMapPropInput{
			Kind:              kind,
			Type:              childTypeRef,
			PrefixForChildren: result.Name,
		})
		if err != nil {
			return nil, err
		}

		result.Props[k] = resolvedProp
	}

	biggestPropName := 0
	biggestPropType := 0
	for _, v := range result.Props {
		if len(v.Name) > biggestPropName {
			biggestPropName = len(v.Name)
		}
		if len(v.Type) > biggestPropType {
			biggestPropType = len(v.Type)
		}
	}

	for _, v := range result.Props {
		v.Spacing1 = strings.Repeat(" ", biggestPropName-len(v.Name))
		v.Spacing2 = strings.Repeat(" ", biggestPropType-len(v.Type))
	}

	if kind == Kind_Entity {
		self.Entities = append(self.Entities, result)
	} else if kind == Kind_Repository {
		self.TypesRepository = append(self.TypesRepository, result)
		self.TypesRepositoryToAvoidDuplication[name] = result
	} else if kind == Kind_Usecase {
		self.TypesUsecase = append(self.TypesUsecase, result)
		self.TypesUsecaseToAvoidDuplication[name] = result
	}

	return result, nil
}
