package types_parser

import (
	"sort"
)

func (self *typeParser) GetEnums() []*Enum {
	sort.Slice(self.enums, func(i, j int) bool {
		return self.enums[i].GolangName < self.enums[j].GolangName
	})
	return self.enums
}

func (self *typeParser) GetTypes() []*Type {
	sort.Slice(self.types, func(i, j int) bool {
		return self.types[i].GolangType < self.types[j].GolangType
	})
	return self.types
}

func (self *typeParser) GetEvents() []*Type {
	sort.Slice(self.events, func(i, j int) bool {
		return self.events[i].GolangType < self.events[j].GolangType
	})
	return self.events
}

func (self *typeParser) GetEntities() []*Type {
	sort.Slice(self.entities, func(i, j int) bool {
		return self.entities[i].GolangType < self.entities[j].GolangType
	})
	return self.entities
}

func (self *typeParser) GetRepository() []*Type {
	sort.Slice(self.repository, func(i, j int) bool {
		return self.repository[i].GolangType < self.repository[j].GolangType
	})
	return self.repository
}

func (self *typeParser) GetUsecase() []*Type {
	sort.Slice(self.usecase, func(i, j int) bool {
		return self.usecase[i].GolangType < self.usecase[j].GolangType
	})
	return self.usecase
}
