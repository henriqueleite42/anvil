package types

type Schema struct {
	Domain     string
	Types      *Types
	Enums      *Enums
	Entities   *Entities
	Repository *Repository
	Events     *Events
	Usecase    *Usecase
	Delivery   *Delivery
}
