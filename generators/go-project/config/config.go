package generator_config

import (
	"log"
	"os"

	"github.com/henriqueleite42/anvil/language-helpers/golang/imports"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
	"gopkg.in/yaml.v3"
)

const GENERATOR_NAME = "go-project"

type GeneratorConfig struct {
	OutDir         *string         `yaml:"OutDir"`
	ProjectName    string          `yaml:"ProjectName"`
	GoVersion      string          `yaml:"GoVersion"`
	PbModuleImport *imports.Import // PbModulePath resolved
}

func GetConfig(filePath string) *GeneratorConfig {
	configFileData, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	config := &schemas.Config{}
	err = yaml.Unmarshal(configFileData, &config)
	if err != nil {
		log.Fatal(err)
	}

	var params map[string]any = nil
	for _, v := range config.Generators {
		if v.Name != GENERATOR_NAME {
			continue
		}

		params = v.Parameters
		break
	}
	if params == nil {
		return &GeneratorConfig{}
	}

	var outDir *string = nil
	outDirAny, ok := params["OutDir"]
	if ok {
		outDirString, ok := outDirAny.(string)
		if !ok {
			log.Fatalf("%s: fail to parse OutDir to string", GENERATOR_NAME)
		}
		outDir = &outDirString
	}

	moduleNameAny, ok := params["ProjectName"]
	if !ok {
		log.Fatalf("%s: ProjectName is a required parameter", GENERATOR_NAME)
	}
	ProjectName, ok := moduleNameAny.(string)
	if !ok {
		log.Fatalf("%s: fail to parse ProjectName to string", GENERATOR_NAME)
	}

	goVersionAny, ok := params["GoVersion"]
	if !ok {
		log.Fatalf("%s: GoVersion is a required parameter", GENERATOR_NAME)
	}
	goVersion, ok := goVersionAny.(string)
	if !ok {
		log.Fatalf("%s: fail to parse GoVersion to string", GENERATOR_NAME)
	}

	var pbModulePath string
	pbModulePathAny, ok := params["PbModulePath"]
	if ok {
		pbModulePathString, ok := pbModulePathAny.(string)
		if !ok {
			log.Fatalf("%s: fail to parse PbModulePath to string", GENERATOR_NAME)
		}
		pbModulePath = pbModulePathString
	} else {
		pbModulePath = ProjectName + "/internal/delivery/grpc/pb"
	}
	pbModuleImport := imports.NewImport(pbModulePath, nil)

	return &GeneratorConfig{
		OutDir:         outDir,
		ProjectName:    ProjectName,
		GoVersion:      goVersion,
		PbModuleImport: pbModuleImport,
	}
}
