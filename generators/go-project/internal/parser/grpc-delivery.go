package parser

import (
	"fmt"
	"sort"
	"strings"

	"github.com/henriqueleite42/anvil/generators/go-project/internal/templates"
	"github.com/henriqueleite42/anvil/language-helpers/golang/formatter"
	"github.com/henriqueleite42/anvil/language-helpers/golang/grpc"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

type GrpcDeliveryParsed struct {
	Imports [][]string
	Methods []*templates.TemplMethodDelivery
}

func (self *Parser) resolveGrpcDelivery(dlv *schemas.DeliveryGrpcRpc, curDomain string) error {
	methodName := dlv.Name

	method, ok := self.Schema.Usecases.Usecases[curDomain].Methods.Methods[dlv.UsecaseMethodHash]
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
		Domain:      curDomain,
		DomainCamel: formatter.PascalToCamel(curDomain),
		DomainSnake: formatter.PascalToSnake(curDomain),
		MethodName:  methodName,
		Input:       input,
		Output:      output,
		Order:       dlv.Order,
	})

	return nil
}

func (self *Parser) ParseDeliveriesGrpc(curDomain string) (*GrpcDeliveryParsed, error) {
	if self.Schema.Deliveries == nil || self.Schema.Deliveries.Deliveries == nil {
		return &GrpcDeliveryParsed{}, nil
	}

	deliveries, ok := self.Schema.Deliveries.Deliveries[curDomain]
	if !ok {
		return &GrpcDeliveryParsed{}, nil
	}

	if deliveries.Grpc == nil {
		return &GrpcDeliveryParsed{}, nil
	}

	self.MethodsGrpcDelivery = []*templates.TemplMethodDelivery{}

	for _, v := range deliveries.Grpc.Rpcs {
		if !strings.HasPrefix(v.Ref, curDomain) {
			continue
		}

		err := self.resolveGrpcDelivery(v, curDomain)
		if err != nil {
			return nil, err
		}
	}

	self.GoTypesParserUsecase.AddImport("context")
	self.GoTypesParserUsecase.AddImport("github.com/rs/xid")
	self.GoTypesParserUsecase.AddImport("github.com/rs/zerolog")
	self.GoTypesParserUsecase.AddImport("google.golang.org/grpc")
	imports := self.GoTypesParserUsecase.GetImports()
	self.GoTypesParserUsecase.ResetImports()

	sort.Slice(self.MethodsGrpcDelivery, func(i, j int) bool {
		return self.MethodsGrpcDelivery[i].Order < self.MethodsGrpcDelivery[j].Order
	})

	return &GrpcDeliveryParsed{
		Imports: imports,
		Methods: self.MethodsGrpcDelivery,
	}, nil
}
