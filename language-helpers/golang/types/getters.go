package types_parser

import "sort"

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

func (self *typeParser) GetTypes() []*Type {
	return self.types
}

func (self *typeParser) GetEvents() []*Type {
	return self.events
}

func (self *typeParser) GetEntities() []*Type {
	return self.entities
}

func (self *typeParser) GetRepository() []*Type {
	return self.repository
}

func (self *typeParser) GetUsecase() []*Type {
	return self.usecase
}
