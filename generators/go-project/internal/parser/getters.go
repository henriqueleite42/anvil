package parser

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ettle/strcase"
	"github.com/henriqueleite42/anvil/generators/go-project/internal/templates"
	types_parser "github.com/henriqueleite42/anvil/language-helpers/golang/types"
)

func typeToTemplType(curPkg string, t *types_parser.Type) *templates.TemplType {
	var mapProps []*templates.TemplTypeMapProp = nil
	if t.MapProps != nil {

		biggestName := 0
		biggestType := 0
		for _, v := range t.MapProps {
			lenName := len(v.Name)
			if lenName > biggestName {
				biggestName = lenName
			}
			tType := v.Type.GetFullTypeName(curPkg)
			lenType := len(tType)
			if lenType > biggestType {
				biggestType = lenType
			}
		}

		mapProps = make([]*templates.TemplTypeMapProp, 0, len(t.MapProps))

		for _, v := range t.MapProps {
			tType := v.Type.GetFullTypeName(curPkg)
			mapProps = append(mapProps, &templates.TemplTypeMapProp{
				Name:     v.Name,
				Spacing1: strings.Repeat(" ", biggestName-len(v.Name)),
				Type:     tType,
				Spacing2: strings.Repeat(" ", biggestType-len(tType)),
				Tags:     v.GetTagsString(),
			})
		}
	}

	return &templates.TemplType{
		AnvilType:          t.AnvilType,
		ModuleImport:       t.ModuleImport,
		GolangType:         t.GolangType,
		Optional:           t.Optional,
		MapProps:           mapProps,
		ImportsUnorganized: t.GetImportsUnorganized(),
	}
}

func enumToTemplEnum(e *types_parser.Enum) *templates.TemplEnum {
	var values []*templates.TemplEnumValue = nil
	if e.Values != nil {
		biggestName := 0
		for _, v := range e.Values {
			name := fmt.Sprintf("%s_%s", e.GolangName, v.Name)
			lenName := len(name)
			if lenName > biggestName {
				biggestName = lenName
			}
		}

		values = make([]*templates.TemplEnumValue, 0, len(e.Values))

		for _, v := range e.Values {
			fullName := fmt.Sprintf("%s_%s", e.GolangName, v.Name)
			values = append(values, &templates.TemplEnumValue{
				Name:    v.Name,
				Spacing: strings.Repeat(" ", biggestName-len(fullName)),
				Idx:     v.Idx,
				Value:   v.Value,
			})
		}
	}

	var deprecatedValues []*templates.TemplEnumValue = nil
	if e.DeprecatedValues != nil {
		biggestName := 0
		for _, v := range e.DeprecatedValues {
			name := fmt.Sprintf("%s_%s", e.GolangName, v.Name)
			lenName := len(name)
			if lenName > biggestName {
				biggestName = lenName
			}
		}

		deprecatedValues = make([]*templates.TemplEnumValue, 0, len(e.DeprecatedValues))

		for _, v := range e.DeprecatedValues {
			deprecatedValues = append(deprecatedValues, &templates.TemplEnumValue{
				Name:    v.Name,
				Spacing: strings.Repeat(" ", biggestName-len(v.Name)),
				Idx:     v.Idx,
				Value:   v.Value,
			})
		}
	}

	return &templates.TemplEnum{
		AnvilEnum: e.AnvilEnum,

		GolangName:       e.GolangName,
		GolangType:       e.GolangType,
		Values:           values,
		DeprecatedValues: deprecatedValues,
	}
}

func (self *Parser) GetEnums() (map[string][]*templates.TemplEnum, error) {
	enumsByDomain := map[string][]*templates.TemplEnum{}
	enums := self.GoTypesParser.GetEnums()
	for _, e := range enums {
		if _, ok := enumsByDomain[e.AnvilEnum.Domain]; !ok {
			enumsByDomain[e.AnvilEnum.Domain] = []*templates.TemplEnum{}
		}

		templEnum := enumToTemplEnum(e)

		enumsByDomain[e.AnvilEnum.Domain] = append(enumsByDomain[e.AnvilEnum.Domain], templEnum)
	}
	for domain := range enumsByDomain {
		sort.Slice(enumsByDomain[domain], func(i, j int) bool {
			return enumsByDomain[domain][i].GolangType < enumsByDomain[domain][j].GolangType
		})
	}

	return enumsByDomain, nil
}

func (self *Parser) GetTypes() (map[string][]*templates.TemplType, error) {
	typesByDomain := map[string][]*templates.TemplType{}
	types := self.GoTypesParser.GetTypes()
	for _, t := range types {
		if _, ok := typesByDomain[t.AnvilType.Domain]; !ok {
			typesByDomain[t.AnvilType.Domain] = []*templates.TemplType{}
		}

		templType := typeToTemplType("models", t)

		typesByDomain[t.AnvilType.Domain] = append(typesByDomain[t.AnvilType.Domain], templType)
	}
	for domain := range typesByDomain {
		sort.Slice(typesByDomain[domain], func(i, j int) bool {
			return typesByDomain[domain][i].GolangType < typesByDomain[domain][j].GolangType
		})
	}

	return typesByDomain, nil
}

func (self *Parser) GetEntities() (map[string][]*templates.TemplType, error) {
	entitiesByDomain := map[string][]*templates.TemplType{}
	entities := self.GoTypesParser.GetEntities()
	for _, t := range entities {
		if _, ok := entitiesByDomain[t.AnvilType.Domain]; !ok {
			entitiesByDomain[t.AnvilType.Domain] = []*templates.TemplType{}
		}

		templType := typeToTemplType("models", t)

		entitiesByDomain[t.AnvilType.Domain] = append(entitiesByDomain[t.AnvilType.Domain], templType)
	}
	for domain := range entitiesByDomain {
		sort.Slice(entitiesByDomain[domain], func(i, j int) bool {
			return entitiesByDomain[domain][i].GolangType < entitiesByDomain[domain][j].GolangType
		})
	}

	return entitiesByDomain, nil
}

func (self *Parser) GetEvents() (map[string][]*templates.TemplType, error) {
	eventsByDomain := map[string][]*templates.TemplType{}
	events := self.GoTypesParser.GetEvents()
	for _, t := range events {
		if _, ok := eventsByDomain[t.AnvilType.Domain]; !ok {
			eventsByDomain[t.AnvilType.Domain] = []*templates.TemplType{}
		}

		templType := typeToTemplType("models", t)

		eventsByDomain[t.AnvilType.Domain] = append(eventsByDomain[t.AnvilType.Domain], templType)
	}
	for domain := range eventsByDomain {
		sort.Slice(eventsByDomain[domain], func(i, j int) bool {
			return eventsByDomain[domain][i].GolangType < eventsByDomain[domain][j].GolangType
		})
	}

	return eventsByDomain, nil
}

func (self *Parser) GetRepositoryTypes() (map[string][]*templates.TemplType, error) {
	repositoryTypesByDomain := map[string][]*templates.TemplType{}
	repositoryTypes := self.GoTypesParser.GetRepository()
	for _, t := range repositoryTypes {
		if _, ok := repositoryTypesByDomain[t.AnvilType.Domain]; !ok {
			repositoryTypesByDomain[t.AnvilType.Domain] = []*templates.TemplType{}
		}

		templType := typeToTemplType(strcase.ToSnake(t.AnvilType.Domain)+"_repository", t)

		repositoryTypesByDomain[t.AnvilType.Domain] = append(repositoryTypesByDomain[t.AnvilType.Domain], templType)
	}
	for domain := range repositoryTypesByDomain {
		sort.Slice(repositoryTypesByDomain[domain], func(i, j int) bool {
			return repositoryTypesByDomain[domain][i].GolangType < repositoryTypesByDomain[domain][j].GolangType
		})
	}

	return repositoryTypesByDomain, nil
}

func (self *Parser) GetUsecaseTypes() (map[string][]*templates.TemplType, error) {
	usecaseTypesByDomain := map[string][]*templates.TemplType{}
	usecaseTypes := self.GoTypesParser.GetUsecase()
	for _, t := range usecaseTypes {
		if _, ok := usecaseTypesByDomain[t.AnvilType.Domain]; !ok {
			usecaseTypesByDomain[t.AnvilType.Domain] = []*templates.TemplType{}
		}

		templType := typeToTemplType(strcase.ToSnake(t.AnvilType.Domain)+"_usecase", t)

		usecaseTypesByDomain[t.AnvilType.Domain] = append(usecaseTypesByDomain[t.AnvilType.Domain], templType)
	}
	for domain := range usecaseTypesByDomain {
		sort.Slice(usecaseTypesByDomain[domain], func(i, j int) bool {
			return usecaseTypesByDomain[domain][i].GolangType < usecaseTypesByDomain[domain][j].GolangType
		})
	}

	return usecaseTypesByDomain, nil
}

func (self *Parser) GetRepositories() (map[string]*ParserRepository, error) {
	return self.repositories, nil
}

func (self *Parser) GetUsecases() (map[string]*ParserUsecase, error) {
	return self.usecases, nil
}

func (self *Parser) GetGrpcDeliveries() (map[string]*ParserGrpcDelivery, error) {
	return self.grpcDeliveries, nil
}

func (self *Parser) GetHttpDeliveries() (map[string]*ParserHttpDelivery, error) {
	return self.httpDeliveries, nil
}

func (self *Parser) GetQueueDeliveries() (map[string]*ParserQueueDelivery, error) {
	return self.queueDeliveries, nil
}
