package generator_config

import (
	"log"
	"os"

	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
	"gopkg.in/yaml.v3"
)

const GENERATOR_NAME = "go-project"

type GeneratorConfig struct {
	OutDir     *string `yaml:"OutDir"`
	ModuleName string  `yaml:"ModuleName"`
	GoVersion  string  `yaml:"GoVersion"`
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

	moduleNameAny, ok := params["ModuleName"]
	if !ok {
		log.Fatalf("%s: ModuleName is a required parameter", GENERATOR_NAME)
	}
	moduleName, ok := moduleNameAny.(string)
	if !ok {
		log.Fatalf("%s: fail to parse ModuleName to string", GENERATOR_NAME)
	}

	goVersionAny, ok := params["GoVersion"]
	if !ok {
		log.Fatalf("%s: GoVersion is a required parameter", GENERATOR_NAME)
	}
	goVersion, ok := goVersionAny.(string)
	if !ok {
		log.Fatalf("%s: fail to parse GoVersion to string", GENERATOR_NAME)
	}

	return &GeneratorConfig{
		OutDir:     outDir,
		ModuleName: moduleName,
		GoVersion:  goVersion,
	}
}
