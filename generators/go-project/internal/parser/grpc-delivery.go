package parser

import (
	"fmt"

	"github.com/henriqueleite42/anvil/generators/go-project/internal/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/formatter"
	"github.com/henriqueleite42/anvil/language-helpers/golang/grpc"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func (self *Parser) ResolveGrpcDelivery(dlv *schemas.DeliveryGrpcRpc) error {
	methodName := dlv.Name

	method, ok := self.Schema.Usecase.Methods.Methods[dlv.UsecaseMethodHash]
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
		Schema:       self.Schema,
		GoTypeParser: self.GoTypesParserUsecase,
	})

	var input *grpc.Type = nil
	if inputTypeHash != "" {
		t, ok := self.Schema.Types.Types[inputTypeHash]
		if !ok {
			return fmt.Errorf("type \"%s\" not found for input of method \"%s\"", inputTypeHash, methodName)
		}

		templT, err := grpcParser.ProtoToGo(&grpc.ProtoToGoInput{
			Type:                     t,
			MethodName:               methodName,
			VariableToAccessTheValue: "i",
			HasOutput:                true,
		})
		if err != nil {
			return err
		}

		input = templT
		self.GoTypesParserUsecase.AddImport("errors")
	} else {
		self.GoTypesParserUsecase.AddImport("google.golang.org/protobuf/types/known/emptypb")
	}

	var output *grpc.Type = nil
	if outputTypeHash != "" {
		t, ok := self.Schema.Types.Types[outputTypeHash]
		if !ok {
			return fmt.Errorf("type \"%s\" not found for output of method \"%s\"", outputTypeHash, methodName)
		}

		templT, err := grpcParser.GoToProto(&grpc.GoToProtoInput{
			Type:                     t,
			MethodName:               methodName,
			VariableToAccessTheValue: "result",
			HasOutput:                true,
		})
		if err != nil {
			return err
		}

		output = templT
	} else {
		self.GoTypesParserUsecase.AddImport("google.golang.org/protobuf/types/known/emptypb")
	}

	self.MethodsGrpcDelivery = append(self.MethodsGrpcDelivery, &templates.TemplMethodDelivery{
		Domain:      self.Schema.Domain,
		DomainCamel: formatter.PascalToCamel(self.Schema.Domain),
		DomainSnake: formatter.PascalToSnake(self.Schema.Domain),
		MethodName:  methodName,
		Input:       input,
		Output:      output,
		Order:       dlv.Order,
	})

	return nil
}
