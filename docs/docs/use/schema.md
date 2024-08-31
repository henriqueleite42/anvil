---
sidebar_position: 10
---

# Schema

## Linter & Autocomplete

The best option of linter for now is [YAML Language Support](https://github.com/redhat-developer/yaml-language-server) by Red Hat.

They have extensions for the most famous code editors:
- Zed: Already installed by default
- VSCode: [Extension](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml)
- Neovim: You must [configure the LSP manually](https://neovim.io/doc/user/lsp.html).

To configure the LSP, you must specify a JSON schema that we provide:
```
https://github.com/henriqueleite42/anvil/cli/blob/master/schemas/v1.0.0.json
```

Heres an example of how to configure it in VSCode and Zed, at the begging of your `.anv` file, add this:
```
# yaml-language-server: $schema=https://github.com/henriqueleite42/anvil/cli/blob/master/schemas/v1.0.0.json
```

:::danger

Be sure to replace the schema version to the latest version.

:::

## Parts

- Domain: Name of the domain being documented
- Relationships: The relationships that your domain have with other domains and micro-services
- Types: Generic types to be used as Input or Output for `Repository` and `Usecase` methods
- Enums: Enums to be used as types for `Entity`, `Events` and `Repository` and `Usecase` methods
- Entities: Tables on your database
- Repository: The way that you communicate with the tables on your database
- Events: Events emitted by this domain
- Usecase: Where all the business logic stays, has the methods to be used by the consumers through delivery methods
- Delivery: General config for delivery methods

## Recommended Confidentiality Levels

- Low: Can be accessed by anyone with access to the service, can be logged and send in events
- Medium: Can only be accessed by services with special permission, cannot be logged or send in events
- High: Can only be accessed inside the domain, cannot be accessed by other services, logged or send in events

## `$ref`

- The only things that you can `$ref` from relationships are `Events` and `Entities`, using the format: `Relationships.Foo.Entities.Bar`. `Foo` must be replaced by the name of the relationship, and **NOT** the `Domain` specified in the relationship file.
