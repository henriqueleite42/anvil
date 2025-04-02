package parser

import (
	"sort"
)

func (self *Parser) sortMethods() {
	if self.repositories != nil {
		for domain := range self.repositories {
			if _, ok := self.repositories[domain]; !ok {
				continue
			}
			if self.repositories[domain].Methods == nil {
				continue
			}

			sort.Slice(self.repositories[domain].Methods, func(i, j int) bool {
				return self.repositories[domain].Methods[i].Order < self.repositories[domain].Methods[j].Order
			})
		}
	}

	if self.usecases != nil {
		for domain := range self.usecases {
			if _, ok := self.usecases[domain]; !ok {
				continue
			}
			if self.usecases[domain].Methods == nil {
				continue
			}

			sort.Slice(self.usecases[domain].Methods, func(i, j int) bool {
				return self.usecases[domain].Methods[i].Order < self.usecases[domain].Methods[j].Order
			})
		}
	}

	if self.grpcDeliveries != nil {
		for domain := range self.grpcDeliveries {
			if _, ok := self.grpcDeliveries[domain]; !ok {
				continue
			}
			if self.grpcDeliveries[domain].Methods == nil {
				continue
			}

			sort.Slice(self.grpcDeliveries[domain].Methods, func(i, j int) bool {
				return self.grpcDeliveries[domain].Methods[i].Order < self.grpcDeliveries[domain].Methods[j].Order
			})
		}
	}
}

func (self *Parser) Parse() error {
	err := self.parseEnums()
	if err != nil {
		return err
	}

	err = self.parseTypes()
	if err != nil {
		return err
	}

	err = self.parseEntities()
	if err != nil {
		return err
	}

	err = self.parseRepositories()
	if err != nil {
		return err
	}

	err = self.parseUsecases()
	if err != nil {
		return err
	}

	err = self.parseDeliveriesGrpc(self.config)
	if err != nil {
		return err
	}

	err = self.parseDeliveriesQueue()
	if err != nil {
		return err
	}

	self.sortMethods()

	return nil
}
