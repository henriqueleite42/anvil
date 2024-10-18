// Credits to:
// https://github.com/cnf/structhash/blob/master/structhash.go
// Unfortunately, it's to old to make it installable in any other projects,
// so I had to copy+paste it

package hashing

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"reflect"
	"sort"
	"strconv"
)

func Struct(c interface{}) (string, error) {
	sum := sha1.Sum(serialize(c))

	return hex.EncodeToString(sum[:]), nil
}

type item struct {
	name  string
	value reflect.Value
}

type itemSorter []item

func (s itemSorter) Len() int {
	return len(s)
}

func (s itemSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s itemSorter) Less(i, j int) bool {
	return s[i].name < s[j].name
}

func writeValue(buf *bytes.Buffer, val reflect.Value) {
	switch val.Kind() {
	case reflect.String:
		buf.WriteByte('"')
		buf.WriteString(val.String())
		buf.WriteByte('"')
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		buf.WriteString(strconv.FormatInt(val.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		buf.WriteString(strconv.FormatUint(val.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		buf.WriteString(strconv.FormatFloat(val.Float(), 'E', -1, 64))
	case reflect.Bool:
		if val.Bool() {
			buf.WriteByte('t')
		} else {
			buf.WriteByte('f')
		}
	case reflect.Ptr:
		if !val.IsNil() || val.Type().Elem().Kind() == reflect.Struct {
			writeValue(buf, reflect.Indirect(val))
		} else {
			writeValue(buf, reflect.Zero(val.Type().Elem()))
		}
	case reflect.Array, reflect.Slice:
		buf.WriteByte('[')
		len := val.Len()
		for i := 0; i < len; i++ {
			if i != 0 {
				buf.WriteByte(',')
			}
			writeValue(buf, val.Index(i))
		}
		buf.WriteByte(']')
	case reflect.Map:
		mk := val.MapKeys()
		items := make([]item, len(mk))
		// Get all values
		for i := range items {
			items[i].name = formatValue(mk[i])
			items[i].value = val.MapIndex(mk[i])
		}

		// Sort values by key
		sort.Sort(itemSorter(items))

		buf.WriteByte('[')
		for i := range items {
			if i != 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(items[i].name)
			buf.WriteByte(':')
			writeValue(buf, items[i].value)
		}
		buf.WriteByte(']')
	case reflect.Struct:
		vtype := val.Type()
		flen := vtype.NumField()
		items := make([]item, 0, flen)
		// Get all fields
		for i := 0; i < flen; i++ {
			field := vtype.Field(i)
			it := item{field.Name, val.Field(i)}
			items = append(items, it)
		}
		// Sort fields by name
		sort.Sort(itemSorter(items))

		buf.WriteByte('{')
		for i := range items {
			if i != 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(items[i].name)
			buf.WriteByte(':')
			writeValue(buf, items[i].value)
		}
		buf.WriteByte('}')
	case reflect.Interface:
		if !val.CanInterface() {
			return
		}
		writeValue(buf, reflect.ValueOf(val.Interface()))
	default:
		buf.WriteString(val.String())
	}
}

func formatValue(val reflect.Value) string {
	if val.Kind() == reflect.String {
		return "\"" + val.String() + "\""
	}

	var buf bytes.Buffer
	writeValue(&buf, val)

	return buf.String()
}

func serialize(object interface{}) []byte {
	var buf bytes.Buffer

	writeValue(
		&buf,
		reflect.ValueOf(object),
	)

	return buf.Bytes()
}
