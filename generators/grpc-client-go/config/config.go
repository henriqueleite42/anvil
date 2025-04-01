package generator_config

import (
	"log"
	"os"

	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
	"gopkg.in/yaml.v3"
)

const GENERATOR_NAME = "grpc-client-go"

type GeneratorConfig struct {
	OutDir      *string `yaml:"OutDir"`
	ProjectName string  `yaml:"ProjectName"`
	ProtoPath   string  `yaml:"ProtoPath"`
	ClientsPath string  `yaml:"ClientsPath"`
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

	projectNameAny, ok := params["ProjectName"]
	if !ok {
		log.Fatalf("%s: ProjectName is a required parameter", GENERATOR_NAME)
	}
	projectName, ok := projectNameAny.(string)
	if !ok {
		log.Fatalf("%s: fail to parse ProjectName to string", GENERATOR_NAME)
	}

	protoPathAny, ok := params["ProtoPath"]
	if !ok {
		log.Fatalf("%s: ProtoPath is a required parameter", GENERATOR_NAME)
	}
	protoPath, ok := protoPathAny.(string)
	if !ok {
		log.Fatalf("%s: fail to parse ProtoPath to string", GENERATOR_NAME)
	}

	clientsPathAny, ok := params["ClientsPath"]
	if !ok {
		log.Fatalf("%s: ClientsPath is a required parameter", GENERATOR_NAME)
	}
	clientsPath, ok := clientsPathAny.(string)
	if !ok {
		log.Fatalf("%s: fail to parse ClientsPath to string", GENERATOR_NAME)
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
		OutDir:      outDir,
		ProjectName: projectName,
		ProtoPath:   protoPath,
		ClientsPath: clientsPath,
	}
}
