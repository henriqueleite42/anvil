package parser

import (
	"fmt"

	generator_config "github.com/henriqueleite42/anvil/generators/go-project/config"
	"github.com/henriqueleite42/anvil/generators/go-project/internal/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/formatter"
	"github.com/henriqueleite42/anvil/language-helpers/golang/grpc"
	"github.com/henriqueleite42/anvil/language-helpers/golang/imports"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func (self *Parser) resolveGrpcDelivery(
	dlv *schemas.DeliveryGrpcRpc,
	curModuleImport *imports.Import,
	pbModuleImport *imports.Import,
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

	var inputTypeHash string
	if method.Input != nil {
		inputTypeHash = method.Input.TypeHash
	}
	var outputTypeHash string
	if method.Output != nil {
		outputTypeHash = method.Output.TypeHash
	}

	grpcParser := grpc.NewGrpcParser(&grpc.NewGrpcParserInput{
		Schema:       self.schema,
		GoTypeParser: self.goTypesParser,
	})

	var input *grpc.ConvertedValue = nil
	if inputTypeHash != "" {
		t, ok := self.schema.Types.Types[inputTypeHash]
		if !ok {
			return fmt.Errorf("type \"%s\" not found for input of method \"%s\"", inputTypeHash, dlv.Name)
		}

		templT, err := grpcParser.ProtoToGo(&grpc.ConverterInput{
			Type:            t,
			CurModuleImport: curModuleImport,
			PbModuleImport:  pbModuleImport,
			VarToConvert:    "i",
		})
		if err != nil {
			return err
		}

		input = templT
	}

	var output *grpc.ConvertedValue = nil
	if outputTypeHash != "" {
		t, ok := self.schema.Types.Types[outputTypeHash]
		if !ok {
			return fmt.Errorf("type \"%s\" not found for output of method \"%s\"", outputTypeHash, dlv.Name)
		}

		templT, err := grpcParser.GoToProto(&grpc.ConverterInput{
			Type:            t,
			CurModuleImport: curModuleImport,
			PbModuleImport:  pbModuleImport,
			VarToConvert:    "result",
		})
		if err != nil {
			return err
		}

		output = templT
	}

	self.grpcDeliveries[dlv.Domain].Methods = append(self.grpcDeliveries[dlv.Domain].Methods, &templates.TemplMethodDelivery{
		Domain:      dlv.Domain,
		DomainCamel: formatter.PascalToCamel(dlv.Domain),
		DomainSnake: formatter.PascalToSnake(dlv.Domain),
		MethodName:  dlv.Name,
		Input:       input,
		Output:      output,
		Order:       dlv.Order,
	})

	return nil
}

func (self *Parser) parseDeliveriesGrpc(config *generator_config.GeneratorConfig) error {
	if self.schema.Deliveries == nil || self.schema.Deliveries.Deliveries == nil {
		return nil
	}

	for _, deliveries := range self.schema.Deliveries.Deliveries {
		if deliveries.Grpc == nil || deliveries.Grpc.Rpcs == nil {
			continue
		}

		for _, v := range deliveries.Grpc.Rpcs {
			if _, ok := self.grpcDeliveries[v.Domain]; !ok {
				self.grpcDeliveries[v.Domain] = &ParserGrpcDelivery{
					Methods: []*templates.TemplMethodDelivery{},
				}
			}

			domainSnake := formatter.PascalToSnake(v.Domain)
			moduleName := domainSnake + "_delivery_grpc"

			curModuleImport := imports.NewImport(config.ModuleName+"/internal/delivery/grpc/"+domainSnake, &moduleName)

			err := self.resolveGrpcDelivery(
				v,
				curModuleImport,
				config.PbModuleImport,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
