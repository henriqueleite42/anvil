package generator_config

import (
	"log"
	"os"

	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
	"gopkg.in/yaml.v3"
)

const GENERATOR_NAME = "grpc-client-go"

type GeneratorConfig struct {
	OutDir            *string `yaml:"OutDir"`
	ModuleName        string  `yaml:"ModuleName"`
	ClientsModuleName string  `yaml:"ClientsModuleName"`
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

	moduleNameAny, ok := params["ModuleName"]
	if !ok {
		log.Fatalf("%s: ModuleName is a required parameter", GENERATOR_NAME)
	}
	moduleName, ok := moduleNameAny.(string)
	if !ok {
		log.Fatalf("%s: fail to parse ModuleName to string", GENERATOR_NAME)
	}

	clientModuleNameAny, ok := params["ClientModuleName"]
	if !ok {
		log.Fatalf("%s: ClientModuleName is a required parameter", GENERATOR_NAME)
	}
	clientModuleName, ok := clientModuleNameAny.(string)
	if !ok {
		log.Fatalf("%s: fail to parse ClientModuleName to string", GENERATOR_NAME)
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

	return &GeneratorConfig{
		OutDir:            outDir,
		ModuleName:        moduleName,
		ClientsModuleName: clientModuleName,
	}
}
