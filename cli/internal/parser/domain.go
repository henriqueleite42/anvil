package parser

import (
	"fmt"

	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func (self *anvToAnvpParser) domain(fileUri string, file map[string]any) (string, error) {
	if self.schema.Schemas == nil {
		self.schema.Schemas = map[string]*schemas.Schema{}
	}

	domainAny, ok := file["Domain"]
	if !ok {
		return "", fmt.Errorf("%s: \"Domain\" must be specified", fileUri)
	}
	domainString, ok := domainAny.(string)
	if !ok {
		return "", fmt.Errorf("%s: fail to parse \"Domain\" to `string`", fileUri)
	}

	var description *string = nil
	descriptionAny, ok := file["Description"]
	if ok {
		descriptionString, ok := descriptionAny.(string)
		if !ok {
			return "", fmt.Errorf("%s: fail to parse \"Description\" to `string`", fileUri)
		}
		description = &descriptionString
	}

	var version *string = nil
	versionAny, ok := file["Version"]
	if ok {
		versionString, ok := versionAny.(string)
		if !ok {
			return "", fmt.Errorf("%s: fail to parse \"Version\" to `string`", fileUri)
		}
		version = &versionString
	}

	self.schema.Schemas[domainString] = &schemas.Schema{
		Domain:      domainString,
		Description: description,
		Version:     version,
		Uri:         fileUri,
	}

	return domainString, nil
}
