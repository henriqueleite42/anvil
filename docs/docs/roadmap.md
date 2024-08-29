---
sidebar_position: 4
---

# Roadmap

## Critical features

Critical and necessary features that are crucial for the project to be used internally.

## Parser

Parse `.anv` files to `.anvp` files(_anvil processed files_). Theses files will be the resolved value of the `.anv` file, with all the standard rules that we apply to the schema.

### Generator gRPC

Generates a proto file for your project.

It replaces the whole protofile, so you can't make changes to it directly.

### Generator Golang Project

Generates a Golang project.

Structure:
```
| internal/
| -- models/
| -- -- implementations/
| -- -- -- example.go
| -- -- example.go
| -- repositories/
| -- -- implementations/
| -- -- -- {domain}/
| -- -- -- -- {domain}.go
| -- -- -- -- {method}.go
| -- -- {domain}.go
| -- usecases/
| -- -- implementations/
| -- -- -- {domain}/
| -- -- -- -- {domain}.go
| -- -- -- -- {method}.go
| -- -- {domain}.go
```

The **models** represents data pof your domain, it includes database tables and events. They are 100% controlled by the generator.

The **repositories** are the way that you interact with your database, it uses the context to get the transaction / db client. The function name, input and output are controlled by the generator, fell free to change everything else.

The **usecases** is were your business logic stays. The function name, input and output are controlled by the generator, fell free to change everything else.

### Generator Golang gRPC delivery

Generates a Golang gRPC delivery.

The delivery is 100% controlled by the generator, so you don't have to write anything, it handles everything for you.

It relies on the usecase pattern created by the [Generator Golang Project](#generator-golang-project).

Structure:
```
| internal/
| -- delivery/
| -- -- {delivery method, like grpc or http}/
| -- -- -- {domain}/
| -- -- -- -- {domain}.go
| -- -- -- -- {route}.go
```

### Generator Golang Pkg gRPC client

Generates a package to make requests to the gRPC API.

It puts all the code in a single file.

Structure:
```
| pkg/
| -- {domain}-api.go
```

### Generator Atlas

Generate [Atlas.io](https://atlasgo.io) schema for postgres

## Extras

Very helpful and necessary features, but that aren't critical for the project internal operation.

### Way to track changes

A git-like way to save changes on the functions and know what has changed.

### Generator changelog with SemVer

Generates a CHANGELOG.md with the changes in the `.anv` files

### Add a more robust authentication

Implement all authentication forms supported by [OpenApi](https://learn.openapis.org/specification/security.html) and more, like AWS Credentials / IAM.

Allow the schema to get properties from the authentication, like the UserId.

### Online platform - v1

- A way to facilitate how you monitor and visualize your system
- Has a list of all your projects
- You can see which project depends on each other
- If we already have CHANGELOGs implemented, see which version of the project it's using

### Three Shaking

Remove unused things from relationships.

Currently, If something has a relationship, we import THE WHOLE FILE and not only the used things. We must find a way to three shake it and generate a smaller file.

### CLI Warnings

- Add warning like "Type Foo is unused, maybe you should remove it."
