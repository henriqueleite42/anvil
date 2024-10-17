package schemas

import "slices"

// Types

type TypeType string

var (
	TypeType_String TypeType = "String"
	TypeType_Bytes  TypeType = "Bytes"

	TypeType_Int   TypeType = "Int"
	TypeType_Int8  TypeType = "Int8"
	TypeType_Int16 TypeType = "Int16"
	TypeType_Int32 TypeType = "Int32"
	TypeType_Int64 TypeType = "Int64"

	TypeType_Uint   TypeType = "Uint"
	TypeType_Uint8  TypeType = "Uint8"
	TypeType_Uint16 TypeType = "Uint16"
	TypeType_Uint32 TypeType = "Uint32"
	TypeType_Uint64 TypeType = "Uint64"

	TypeType_Float   TypeType = "Float"
	TypeType_Float32 TypeType = "Float32"
	TypeType_Float64 TypeType = "Float64"

	TypeType_Bool TypeType = "Bool"

	TypeType_Timestamp TypeType = "Timestamp"

	TypeType_Enum TypeType = "Enum"

	TypeType_Map TypeType = "Map"

	TypeType_List TypeType = "List"
)

var TypeTypeArr = []TypeType{
	TypeType_String,
	TypeType_Bytes,

	TypeType_Int,
	TypeType_Int8,
	TypeType_Int16,
	TypeType_Int32,
	TypeType_Int64,

	TypeType_Uint,
	TypeType_Uint8,
	TypeType_Uint16,
	TypeType_Uint32,
	TypeType_Uint64,

	TypeType_Float,
	TypeType_Float32,
	TypeType_Float64,

	TypeType_Bool,
	TypeType_Timestamp,
	TypeType_Enum,
	TypeType_Map,
	TypeType_List,
}

func ToTypeType(i string) (TypeType, bool) {
	ft := TypeType(i)

	return ft, slices.Contains(TypeTypeArr, ft)
}

type TypeConfidentiality string

var (
	TypeConfidentiality_Low    TypeConfidentiality = "LOW"
	TypeConfidentiality_Medium TypeConfidentiality = "MEDIUM"
	TypeConfidentiality_High   TypeConfidentiality = "HIGH"
)

var TypeConfidentialityArr = []TypeConfidentiality{
	TypeConfidentiality_Low,
	TypeConfidentiality_Medium,
	TypeConfidentiality_High,
}

func ToTypeConfidentiality(i string) (TypeConfidentiality, bool) {
	ft := TypeConfidentiality(i)

	return ft, slices.Contains(TypeConfidentialityArr, ft)
}

// Enums

type EnumType string

var (
	EnumType_String EnumType = "String"

	EnumType_Int   EnumType = "Int"
	EnumType_Int8  EnumType = "Int8"
	EnumType_Int16 EnumType = "Int16"
	EnumType_Int32 EnumType = "Int32"
	EnumType_Int64 EnumType = "Int64"

	EnumType_Uint   EnumType = "Uint"
	EnumType_Uint8  EnumType = "Uint8"
	EnumType_Uint16 EnumType = "Uint16"
	EnumType_Uint32 EnumType = "Uint32"
	EnumType_Uint64 EnumType = "Uint64"
)

// Entities

type NamingCase string

var (
	NamingCase_Snake  NamingCase = "snake"
	NamingCase_Pascal NamingCase = "pascal"
	NamingCase_Camel  NamingCase = "camel"
)

var NamingCaseArr = []NamingCase{
	NamingCase_Snake,
	NamingCase_Pascal,
	NamingCase_Camel,
}

func ToNamingCase(i string) (NamingCase, bool) {
	ft := NamingCase(i)

	return ft, slices.Contains(NamingCaseArr, ft)
}
