package parser

import (
	"fmt"
	"net/url"
	"strings"
	"unicode"

	"github.com/ettle/strcase"
	generator_config "github.com/henriqueleite42/anvil/generators/go-project/config"
	"github.com/henriqueleite42/anvil/generators/go-project/internal/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func pathToPascal(path string) string {
	parsedURL, err := url.Parse(path)
	if err != nil {
		panic(err)
	}

	pathParts := strings.Split(parsedURL.Path, "/")

	pathPascal := ""

	for _, part := range pathParts {
		if part == "" {
			continue
		}

		runes := []rune(strings.ToLower(part))
		runes[0] = unicode.ToUpper(runes[0])
		pathPascal += string(runes)
	}

	return pathPascal
}

func (self *Parser) resolveHttpDelivery(
	dlv *schemas.DeliveryHttpRoute,
	config *generator_config.GeneratorConfig,
) error {
	if self.schema.Usecases == nil ||
		self.schema.Usecases.Usecases == nil {
		return nil
	}
	if _, ok := self.schema.Usecases.Usecases[dlv.Domain]; !ok {
		return nil
	}
	if self.schema.Usecases.Usecases[dlv.Domain].Methods == nil ||
		self.schema.Usecases.Usecases[dlv.Domain].Methods.Methods == nil {
		return nil
	}
	_, ok := self.schema.Usecases.Usecases[dlv.Domain].Methods.Methods[dlv.UsecaseMethodHash]
	if !ok {
		return fmt.Errorf("usecase method \"%s\" not found", dlv.UsecaseMethodHash)
	}

	pathPascal := pathToPascal(dlv.Path)

	self.httpDeliveries[dlv.Domain].Methods = append(self.httpDeliveries[dlv.Domain].Methods, &templates.TemplHttpMethodDelivery{
		DomainPascal:    dlv.Domain,
		DomainCamel:     strcase.ToCamel(dlv.Domain),
		DomainSnake:     strcase.ToSnake(dlv.Domain),
		RouteNamePascal: pathPascal,
		RouteNameSnake:  strcase.ToSnake(pathPascal),
		Order:           dlv.Order,
	})

	return nil
}

func (self *Parser) parseDeliveriesHttp(config *generator_config.GeneratorConfig) error {
	if self.schema.Deliveries == nil || self.schema.Deliveries.Deliveries == nil {
		return nil
	}

	for _, deliveries := range self.schema.Deliveries.Deliveries {
		if deliveries.Http == nil || deliveries.Http.Routes == nil {
			continue
		}

		for _, v := range deliveries.Http.Routes {
			if _, ok := self.httpDeliveries[v.Domain]; !ok {
				self.httpDeliveries[v.Domain] = &ParserHttpDelivery{
					Methods: []*templates.TemplHttpMethodDelivery{},
				}
			}

			err := self.resolveHttpDelivery(
				v,
				config,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
