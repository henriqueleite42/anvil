package types_test

import (
	"fmt"
	"log"

	"github.com/henriqueleite42/anvil/language-helpers/golang/formatter"
	"github.com/henriqueleite42/anvil/language-helpers/golang/hashing"
	"github.com/henriqueleite42/anvil/language-helpers/golang/imports"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
	types_parser "github.com/henriqueleite42/anvil/language-helpers/golang/types"
)

func createBaseInstance(anvp *schemas.AnvpSchema) types_parser.TypesParser {
	tp, err := types_parser.NewTypeParser(&types_parser.NewTypeParserInput{
		Schema: anvp,

		GetEnumsImport: func(e *schemas.Enum) *imports.Import {
			path := fmt.Sprintf("github.com/foo/bar/enums/%s", formatter.PascalToSnake(e.Domain))
			return imports.NewImport(path, nil)
		},
		GetTypesImport: func(t *schemas.Type) *imports.Import {
			path := fmt.Sprintf("github.com/foo/bar/types/%s", formatter.PascalToSnake(t.Domain))
			return imports.NewImport(path, nil)
		},
		GetEventsImport: func(t *schemas.Type) *imports.Import {
			path := fmt.Sprintf("github.com/foo/bar/events/%s", formatter.PascalToSnake(t.Domain))
			return imports.NewImport(path, nil)
		},
		GetEntitiesImport: func(t *schemas.Type) *imports.Import {
			path := fmt.Sprintf("github.com/foo/bar/entities/%s", formatter.PascalToSnake(t.Domain))
			return imports.NewImport(path, nil)
		},
		GetRepositoryImport: func(t *schemas.Type) *imports.Import {
			alias := formatter.PascalToSnake(t.Domain) + "_repository"
			path := fmt.Sprintf("github.com/foo/bar/repository/%s", alias)
			return imports.NewImport(path, nil)
		},
		GetUsecaseImport: func(t *schemas.Type) *imports.Import {
			alias := formatter.PascalToSnake(t.Domain) + "_usecase"
			path := fmt.Sprintf("github.com/foo/bar/usecase/%s", alias)
			return imports.NewImport(path, nil)
		},
	})
	if err != nil {
		log.Fatal("fail to create instance of TypeParser")
	}

	return tp
}

func strOnly() {
	/** Based on schema:
	Domain: Foo

	Types:
		Str:
			Type: String
	*/
	str := &schemas.Type{
		Ref:             "Foo.Types.Str",
		OriginalPath:    "Foo.Types.Str",
		RootNode:        "Types",
		Domain:          "Foo",
		StateHash:       "--",
		Name:            "Str",
		Confidentiality: "LOW",
		Optional:        false,
		Type:            schemas.TypeType_String,
	}

	strRefHash := hashing.String(str.Ref)

	anvp := &schemas.AnvpSchema{
		Schemas: map[string]*schemas.Schema{
			"Foo": {
				Domain: "Foo",
				Uri:    "--",
			},
		},
		Types: &schemas.Types{
			StateHash: "--",
			Types: map[string]*schemas.Type{
				strRefHash: str,
			},
		},
	}

	tp := createBaseInstance(anvp)

	strParsed, err := tp.ParseType(str)
	if err != nil {
		log.Fatalf("TypeParser: fail to parse \"str\": %s", err.Error())
	}

	if strParsed.AnvilType != str.Type {
		log.Fatal("TypeParser: fail to parse \"str\": AnvilType != str.Type")
	}
	if strParsed.GolangType != "string" {
		log.Fatal("TypeParser: fail to parse \"str\": GolangType != string")
	}
	if strParsed.MapProps != nil {
		log.Fatal("TypeParser: fail to parse \"str\": MapProps != nil")
	}
	if strParsed.ModuleImport != nil {
		log.Fatal("TypeParser: fail to parse \"str\": ModuleImport != nil")
	}
	if strParsed.Optional != str.Optional {
		log.Fatal("TypeParser: fail to parse \"str\": Optional != str.Optional")
	}
}

func ParseType() {
	strOnly()
}
