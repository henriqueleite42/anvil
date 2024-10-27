package parser

func (self *Parser) parseEnums() error {
	if self.schema.Enums == nil || self.schema.Enums.Enums == nil {
		return nil
	}

	for _, v := range self.schema.Enums.Enums {
		_, err := self.goTypesParser.ParseEnum(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *Parser) parseTypes() error {
	if self.schema.Types == nil || self.schema.Types.Types == nil {
		return nil
	}

	for _, v := range self.schema.Types.Types {
		_, err := self.goTypesParser.ParseType(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *Parser) parseEntities() error {
	if self.schema.Entities == nil || self.schema.Entities.Entities == nil {
		return nil
	}

	for _, v := range self.schema.Entities.Entities {
		entity := self.schema.Types.Types[v.TypeHash]

		_, err := self.goTypesParser.ParseType(entity)
		if err != nil {
			return err
		}
	}

	return nil
}
