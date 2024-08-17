# Anvil

Anvil is like OpenAPI schemas for microservices, but instead of only documenting your http routes, it helps you to manage most of the aspects of all your microservices at a global scale. It's created to medium~big companies and suffer on delivering things with velocity and consistency.

It follows an _schema-first_ approach, of instead of writing your code first, you write an schema, and it generates most of the code for you, all that parts that are repetitive and doesn't influence in the performance, while letting you have 100% control of the part that matters: the business logic.

The schema is designed for **Event-Oriented, Domain-Driven, Decoupled MicroServices, with a Delivery-Usecase-Repository architecture and SQL Databases**. It still can be used for monoliths and other types of architectures, but we don't maintain the schema to be extremely flexible and a silver bullet for all the projects. Or goal here **IS NOT** to allow creativity, is to have a way to create scalable, secure and maintainable applications.

## How Anvil can help you

- Documentation tool:
  - Were data is used across your whole system and projects
  - If data is confidential / has any legal protection
  - Which services and domains depends on each other
  - Which events and routes does a service has, so you can subscribe to them without having to communicate with the owner team, making process faster
- Code generation tool:
  - Define your own project pattern and generate projects on your own way
  - Ensure that developers follow a specific pattern
  - Keep dependencies updated across all your projects
  - Ensure that best practices / your practices are begin followed
  - Automatic integrate observability tools, logging, anything you want, on your projects, by default
  - Ensure standards on variable names, event names and patterns, folder structure and project architecture
- Refactoring tool:
  - Want to refactor an old project in a new language / pattern? Use the same schema, a different generator, and you only have to copy-paste / do small adjusts on the business logic. Decrease the refactoring time by an immensurable amount of time.
  - Test the same project on different languages and architectures, to see which one is the best. Use one schema, different generators, and generate the same API in multiple languages in a fraction of the time that it would took.

Anvil allows you to write once, document and generate everywhere. Once that you have your schema defined, a lot of doors opens to you.

## What INS'T Anvil

- Something to control/create/update your infrastructure like CloudFormation, Terraform or Serverless Framework
- A framework to magically implement things under the hood, hide complexity and make you dependent on it
- Something to guide exactly how you should implement your code, your architecture, your folder structure, and so on
- A message bus to help you send and receive events

## Why use Anvil

**TL;DR**
Anvil will help you to:
- Need less developers to accomplish the same (probably even best) results
- Decrease the amount of time that it takes to create new products and features, without having to compromise the quality and security of the software
- Better divide the responsibilities of your team, to get the best that they can offer and not needing so many experienced developers to create amazing products

In large organizations, we usually have hundreds or even thousands of micro-services, teams, events, packages and team members changing teams in a daily basis. It's very hard and demanding to maintain everything, to share these knowledge of the best practices, to ensure that all developers not only know how to implement certain patterns but know the way that the company implements certain patterns.

Anvil is created for these kind of ecosystems. It allows you to have one centralized small team of extremely capable developers that say how the things will work, define rules, best practices, standard libraries, and everything else that you need or want to defined, and all the other teams and members of your organization will follow these rules and patterns.

## How Anvil does these things?

Anvil by itself:
- _Schema-first_ approaches help you to visualize the current state of your system in a very easy and fast way: Instead of having to understand code, the project pattern, searching in a bunch of files, go directly to the ONE file definition anf figure it out right away.
- Instead of trying to reinvent the wheel, you can follow a standardized architecture that is scalable, clean, flexible, follows the best practices and allows the work to be divided in multiple steps that can be executed in parallel
- Standardize all you micro-services to follow the exact same patterns for EVERYTHING, decreasing a lot the learning curve and the effort necessary to maintain them
- Allows tracking and usage of confidential and private data, like user's emails, to complain with regulations

Anvil generators:
- Generate e2e tests, useful for early stage startups that can't afford a QA or have enough time to implement more complex tests
- Generate `.proto` files for gRPC APIs and OpenAPI specs for REST APIs
- Generate database migrations and automatically handle them, begin able to easily see the current state of your database by looking at the schema file
- Generate standardized clients for your APIs, with automatically generate `CHANGELOG.md`s that follows [SemVer](https://semver.org)
- You are not stuck to Anvil. It's not a framework, it generates code that you have 100% control of. If you don't want to use Anvil anymore, it will have no impact on your systems.

Anvil plugins:
- Integration with other tools like Jira, Linear, Slack or your own custom system to send notifications / perform tasks
- Easy to change things at a global scale: Do not get stuck into a language, framework or architecture anymore, if you want to change it one day, it will be 100x times easier and faster.

## F.A.Q.

[See the F.A.Q. here.](./FAQ.md)

## Examples

[See the examples here.](./examples.md)

## How Anvil work

Anvil has 4 main parts, each one responsible for a specific complementary role.

### `*.anv` files

The schema definition is a `.anv` file that describes a domain of your service. Each project (micro-service) can have multiple domains in it, and they can be related or not (ideally, if they are in the same project, they should be).

Think about the `.anv` files like a `schema.prisma` or an OpenApi spec, and from this we generate an infinity of things.

### `.anvilconfig`

`.anvilconfig` is the configuration file for Anvil, where you put information like the plugins that you are using, the things that you want to generate, and any other configuration that Anvil CLI or the plugins may need.

It is written in [TOML](https://toml.io/en/).

### CLI

The CLI is the way that you interact with all Anvil things. You can use it to validate your files, generate things, install plugins, run your migrations, and much more.

It's designed to work with CI/CD too 🙌

### Generators

Generators allows you to generate code based on a `.anv` config. They come in various shapes and sizes, and can be used for practically anything:
- Generate an microservice with a specific code pattern, that uses a specific set of libraries
- Generate e2e tests
- Generate changelogs

Generator are were the magic oh Anvil happens.

## Meaning

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
