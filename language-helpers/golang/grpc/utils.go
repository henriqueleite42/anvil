package grpc

import "github.com/henriqueleite42/anvil/language-helpers/golang/schemas"

func isBasicType(t schemas.TypeType) bool {
	return t == schemas.TypeType_String ||
		t == schemas.TypeType_Int ||
		t == schemas.TypeType_Int8 ||
		t == schemas.TypeType_Int16 ||
		t == schemas.TypeType_Int32 ||
		t == schemas.TypeType_Int64 ||
		t == schemas.TypeType_Uint ||
		t == schemas.TypeType_Uint8 ||
		t == schemas.TypeType_Uint16 ||
		t == schemas.TypeType_Uint32 ||
		t == schemas.TypeType_Uint64 ||
		t == schemas.TypeType_Float ||
		t == schemas.TypeType_Float32 ||
		t == schemas.TypeType_Float64 ||
		t == schemas.TypeType_Bool
}