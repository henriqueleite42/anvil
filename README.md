# Hephaestus

Hephaestus is a CLI tool and Schema standard to facilitate the creation and maintainability of micro-services. It's goal is to allow the creation and maintenance of **Event-Oriented, Domain-Driven, Decoupled MicroServices, with a Delivery-Usecase-Repository architecture and SQL Databases**. We can extend it's usability as long as it keeps it's original purpose and maintain the schema as simple as possible.

## Why use Hephaestus

- Has LSP integration and is well documented, easy to understand and get things rolling
- Instead of trying to reinvent the wheel, you can follow a standardized architecture that is scalable, clean, flexible and allows the work to be divided in multiple steps that can be executed in parallel
- Standardize all you micro-services to follow the exact same patterns for EVERYTHING, decreasing a lot the learning curve and the effort necessary to maintain them
- Allows tracking and usage of confidential and private data, like user's emails, to complain with regulations
- Automatically generates e2e tests, the most important tests, extremely useful for early stage startups that can't afford a QA or have enough time to implement more complex tests
- Automatically generates `.proto` files for gRPC APIs and OpenAPI specs for REST APIs
- Generate database migrations and automatically handle them, begin able to easily see the current state of your database by looking at the schema file
- Generate standardized clients for your APIs, with automatically generate `CHANGELOG.md`s that follows [SemVer](https://semver.org)
- Allows for external plugins that allow for integration with other tools like Jira, Linear, Slack or your own custom system

## How Hephaestus work

Hephaestus has 4 main parts, each one responsible for a specific complementary role.

### `.hpt` files

The schema definition is a `.hpt` file that describes a domain of your service. Each project (micro-service) can have multiple domains in it, and they can be related or not (ideally, they should be).

Think about the `.hpt` files like a `schema.prisma` or an OpenApi spec, and from this we generate an infinity of things.

### `.hephaestusconfig`

`.hephaestusconfig` is the configuration file for Hephaestus, where you put informations like the plugins that you are using, the things that you want to generate, and any other configuration that Hephaestus CLI or the plugins may need.

It is written in [TOML](https://toml.io/en/).

### CLI

The CLI is the way that you interact with all Hephaestus things. You can use it to validate your files, generate things, install plugins, run your migrations, and much more.

It's designed to work with CI/CD too 🙌

### Plugins

Plugins allows you to customize Hephaestus as you want. We create the standard, you customize it to attend your needs.

We created some plugins that can be used by anyone independent of the project pattern that they are using, and some specific plugins that fit our needs, but that you can use too.

## F.A.Q.

### Is Hephaestus a silver bullet for every project?

No. We are very clear that Hephaestus is focused on **Event-Oriented, Domain-Driven, Decoupled MicroServices, with a Delivery-Usecase-Repository architecture and SQL Databases**. It may be used in other cases, but we don't give support for these other cases.

### Can Hephaestus be used for monoliths?

Sure, it probably will work great with monoliths too, since they kinda are "big micro-services", we don't guarantee that it will be perfect, but for sure it will help.

### Can Hephaestus be used with NoSql databases?

You kinda can if you have the right plugin, but the schema is not and will not be designed for the specific needs that a NoSql database have. If you want to use a NoSql database as a SQL database, Hephaestus will probably fit your needs.

## Meaning

- Domain: Name of the domain being documented
- Types: Generic types to be used as Input ou Output for `Repository` and `Usecase` methods
- Enums: Enums to be used as types for `Entity`, `Events` and `Repository` and `Usecase` methods
- Entities: Tables on your database
- Repository: The way that you communicate with the tables on your database
- Events: Events emitted by this domain
- Usecase: Where all the business logic stays, has the methods to be used by the consumers through delivery methods
- Delivery: General config for delivery methods

## Rules

- The `Schema` isn't a reflection of the config files, it's the **resolved** version of the config files, so all the functions (like `$ref`) aren't present in the Schema, only the true value of that reference

## Confidentiality Levels

- Low: Can be accessed by anyone with access to the service, can be logged and send in events
- Medium: Can only be accessed by services with special permission, cannot be logged or send in events
- High: Can only be accessed inside the domain, cannot be accessed by other services, logged or send in events
