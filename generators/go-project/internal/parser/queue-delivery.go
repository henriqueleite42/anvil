package parser

import (
	"fmt"
	"strings"

	"github.com/ettle/strcase"
	"github.com/henriqueleite42/anvil/generators/go-project/internal/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
	types_parser "github.com/henriqueleite42/anvil/language-helpers/golang/types"
)

func (self *Parser) resolveQueueDelivery(
	dlv *schemas.DeliveryQueueQueue,
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

	var input *types_parser.Type = nil
	if method.Input != nil {
		t, ok := self.schema.Types.Types[method.Input.TypeHash]
		if !ok {
			return fmt.Errorf("type \"%s\" not found for input of method \"%s\"", method.Input.TypeHash, dlv.QueueId)
		}

		templT, err := self.GoTypesParser.ParseType(t)
		if err != nil {
			return err
		}

		input = templT
	}

	domainCamel := strcase.ToCamel(dlv.Domain)
	inputName := input.GetTypeName("queue")
	// TODO Fix the line bellow. This is an workaround and a terrible way to solve this, but it's fast enough
	inputName = strings.Replace(inputName, "[]*", "", 1)
	self.queueDeliveries[dlv.Domain].Methods = append(self.queueDeliveries[dlv.Domain].Methods, &templates.TemplQueueMethodDelivery{
		DomainPascal:  dlv.Domain,
		DomainCamel:   domainCamel,
		MethodName:    method.Name,
		QueueIdPascal: strcase.ToPascal(dlv.QueueId),
		Input:         input,
		InputName:     inputName,
	})

	return nil
}

func (self *Parser) parseDeliveriesQueue() error {
	if self.schema.Deliveries == nil || self.schema.Deliveries.Deliveries == nil {
		return nil
	}

	for _, deliveries := range self.schema.Deliveries.Deliveries {
		if deliveries.Queue == nil || deliveries.Queue.Queues == nil {
			continue
		}

		for _, v := range deliveries.Queue.Queues {
			if _, ok := self.queueDeliveries[v.Domain]; !ok {
				self.queueDeliveries[v.Domain] = &ParserQueueDelivery{
					Methods: []*templates.TemplQueueMethodDelivery{},
				}
			}

			err := self.resolveQueueDelivery(
				v,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
