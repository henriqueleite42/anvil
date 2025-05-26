package types_parser

import (
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
)

func IsTypePointer(t *schemas.Type) bool {
	if t.Type == schemas.TypeType_Timestamp ||
		t.Type == schemas.TypeType_Enum ||
		t.Type == schemas.TypeType_Map ||
		t.Type == schemas.TypeType_List {
		return true
	}

	return false
}

func IsTypeBool(t *schemas.Type) bool {
	return t.Type == schemas.TypeType_Bool
}
