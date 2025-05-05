package parser

import (
	"fmt"

	"github.com/ettle/strcase"
	generator_config "github.com/henriqueleite42/anvil/generators/go-project/config"
	"github.com/henriqueleite42/anvil/generators/go-project/internal/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

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
	method, ok := self.schema.Usecases.Usecases[dlv.Domain].Methods.Methods[dlv.UsecaseMethodHash]
	if !ok {
		return fmt.Errorf("usecase method \"%s\" not found", dlv.UsecaseMethodHash)
	}

	self.httpDeliveries[dlv.Domain].Methods = append(self.httpDeliveries[dlv.Domain].Methods, &templates.TemplHttpMethodDelivery{
		DomainPascal:    dlv.Domain,
		DomainCamel:     strcase.ToCamel(dlv.Domain),
		DomainSnake:     strcase.ToSnake(dlv.Domain),
		RouteNamePascal: method.Name,
		RouteNameSnake:  strcase.ToSnake(method.Name),
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
