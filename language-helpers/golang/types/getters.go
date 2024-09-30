package types_parser

import "sort"

func (self *typeParser) GetMapTypes() []*Type {
	return self.types
}

func (self *typeParser) GetEnums() []*Enum {
	enums := make([]*Enum, 0, len(self.enums))
	for _, v := range self.enums {
		enums = append(enums, v)
	}
	sort.Slice(enums, func(i, j int) bool {
		return enums[i].GolangType < enums[j].GolangType
	})
	return enums
}
