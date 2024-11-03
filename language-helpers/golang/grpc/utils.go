package grpc

import (
	"fmt"

	"github.com/henriqueleite42/anvil/language-helpers/golang/imports"
	"github.com/henriqueleite42/anvil/language-helpers/golang/schemas"
	types_parser "github.com/henriqueleite42/anvil/language-helpers/golang/types"
)

type EnumConversionType int

const (
	EnumConversionType_FromPb EnumConversionType = iota
	EnumConversionType_ToPb
)

func doestNeedConversion(t schemas.TypeType) bool {
	return t == schemas.TypeType_Bytes ||
		t == schemas.TypeType_String ||
		// t == schemas.TypeType_Int ||
		// t == schemas.TypeType_Int8 ||  // Needs to convert to/from int32
		// t == schemas.TypeType_Int16 ||
		t == schemas.TypeType_Int32 ||
		t == schemas.TypeType_Int64 ||
		// t == schemas.TypeType_Uint ||
		// t == schemas.TypeType_Uint8 ||  // Needs to convert to/from uint32
		// t == schemas.TypeType_Uint16 ||
		t == schemas.TypeType_Uint32 ||
		t == schemas.TypeType_Uint64 ||
		// t == schemas.TypeType_Float || // Needs to convert to/from float32
		t == schemas.TypeType_Float32 ||
		t == schemas.TypeType_Float64 ||
		t == schemas.TypeType_Bool
}

func (self *goGrpcParser) getEnumConvertFuncName(
	way EnumConversionType,
	curModuleImport *imports.Import,
	e *types_parser.Enum,
	varToConvert string,
) string {
	var firstLetter string
	var pkg string

	if self.enumConversionImport.Alias == curModuleImport.Alias {
		firstLetter = "c"
	} else {
		firstLetter = "C"
		pkg = self.enumConversionImport.Alias
	}

	if pkg == "" {
		pkg = pkg + "."
	}

	return fmt.Sprintf("%s%sonvert%sToPb(*%s)", pkg, firstLetter, e.GolangName, varToConvert)
}
