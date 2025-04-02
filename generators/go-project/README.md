# Anvil Generator: GoProject

## Parameters

- **OutDir:** Folder to generate files.
  - Optional, will use current folder if none specified
- **ProjectName:** Golang module name to be used in `go.mod` and imports
- **GoVersion:** Golang version to be used in `go.mod`

## Technical debits

- Generators
	- go-project
		- Queues
			- It considers all queues as `Bulk: true` (can't handle `Bulk: false`)
			- All queues MUST receive an object, it can't work with any other type of input (or lack of)
				- The CLI doesn't warns about this nor the generator
			- It doesn't automatically import / adds as inputs the usecases
			- It doesn't automatically adds the routes to `Listen` method
