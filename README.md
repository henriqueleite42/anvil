# Anvil

Anvil is like OpenAPI schemas for microservices, but instead of only documenting your http routes, it helps you to manage most of the aspects of all your microservices at a global scale. It's created to medium~big companies and suffer on delivering things with velocity and consistency.

It follows an _schema-first_ approach, of instead of writing your code first, you write an schema, and it generates most of the code for you, all that parts that are repetitive and doesn't influence in the performance, while letting you have 100% control of the part that matters: the business logic.

The schema is designed for **Event-Oriented, Domain-Driven, Decoupled MicroServices, with a Delivery-Usecase-Repository architecture and SQL Databases**. It still can be used for monoliths and other types of architectures, but we don't maintain the schema to be extremely flexible and a silver bullet for all the projects. Or goal here **IS NOT** to allow creativity, is to have a way to create scalable, secure and maintainable applications.

## How Anvil can help you

- Documentation tool:
  - Were data is used across your whole system and projects
  - If data is condidential / has any legal protection
  - Which services and domains depends on each other
  - Which events and routes does a service has, so you can subiscrive to them without having to communicate with the owner team, making process faster
- Code generation tool:
  - Define your own project pattern and generate projects on your own way
  - Ensure that developers follow a specific pattern
  - Keep dependencies updated across all your projects
  - Ensure that best practicies / your practices are beign followed
  - Automatic integrate observability tools, logging, anything you want, on your projects, by default
  - Ensure standards on variable names, event names and patterns, folder structe and project architecture
- Refactoring tool:
  - Want to refactor an old project in a new language / pattern? Use the same schema, a different generator, and you only have to copy-paste / do smll adjusts on the business logic. Decrease the refactoring time by an imensurable amount of time.
  - Test the same project on different languages and architectures, to see which one is the best. Use one schema, different generators, and generate the same API in multiple languages in a fraction of the time that it would took.

Anvil allows you to write once, document and generate everywere. Once that you have your schema defined, a lot of doors opens to you.

## What INS'T Anvil

- Something to control/create/update your infrastructure like CloudFormation, Terraform or Serverless Framework
- A framework to magically implement things under the hood, hide complexity and make you dependent on it
- Something to guide exactly how you should implement your code, your architecture, your folder structure, and so on
- A message bus to help you send and receive events

## Why use Anvil

In large organizations, we usually have hundreds or even thousands of micro-services, teams, events, packages and team members changing teams in a daily basis. It's very hard and demanding to maintain everything, to share these knowledge of the best practices, to ensure that all developers not only know how to implement certain patterns but know the way that the company implements certain patterns.

Anvil is created for these kind of ecosystems. It allows you to have one centralized small team of extremely capable developers that say how the things will work, define rules, best practices, standard libraries, and everything else that you need or want to defined, and all the other teams and members of your organization will follow these rules and patterns.

**TL;DR**
Anvil will help you to:
- Need less developers to accomplish the same (probably even best) results
- Decrease the amount of time that it takes to create new products and features, without having to compromise the quality and security of the software
- Better divide the responsibilities of your team, to get the best that they can offer and not needing so many experienced developers to create amazing products

## How Anvil does these things?

- _Schema-first_ approaches help you to visualize the current state of your system in a very easy and fast way: Instead of having to understand code, the project pattern, searching in a bunch of files, go directly to the ONE file definition anf figure it out right away.
- Instead of trying to reinvent the wheel, you can follow a standardized architecture that is scalable, clean, flexible, follows the best practices and allows the work to be divided in multiple steps that can be executed in parallel
- Standardize all you micro-services to follow the exact same patterns for EVERYTHING, decreasing a lot the learning curve and the effort necessary to maintain them
- Allows tracking and usage of confidential and private data, like user's emails, to complain with regulations
- Automatically generates e2e tests, the most important tests, extremely useful for early stage startups that can't afford a QA or have enough time to implement more complex tests
- Automatically generates `.proto` files for gRPC APIs and OpenAPI specs for REST APIs
- Generate database migrations and automatically handle them, begin able to easily see the current state of your database by looking at the schema file
- Generate standardized clients for your APIs, with automatically generate `CHANGELOG.md`s that follows [SemVer](https://semver.org)
- Allows for external plugins that allow for integration with other tools like Jira, Linear, Slack or your own custom system
- Easy to change things at a global scale: Do not get stuck into a language or architecture anymore, if you want to change one day, it's 100x times easier and faster.
- You are not stuck to Anvil. It's not a framework, it generates code that you have 100% control of. If you don't want to use Anvil anymore, it will have no impact on your systems.

## F.A.Q.

### Is Anvil ready to use right out the box?

If you are whiling to adopt one of the architectures proposed by one of ht existing generators, including the ones that we maintain, **yes!**

If you want to implement your own architecture, with your own code patterns, so you will have to write your own generator, what is veeeery complex and time consuming compared to write a single micro-service, but pays off on the long run if you plan to have more than a dozen of micro-services.

### Is Anvil a silver bullet for every project?

No. We are very clear that Anvil is focused on **Event-Oriented, Domain-Driven, Decoupled MicroServices, with a Delivery-Usecase-Repository architecture and SQL Databases**. It may be used in other cases, but we don't give support for these other cases.

### Can Anvil be used for monoliths?

Sure, it probably will work great with monoliths too, since they kinda are "big micro-services", we don't guarantee that it will be perfect, but for sure it will help.

### Can Anvil be used with NoSql databases?

You kinda can if you have the right plugin, but the schema is not and will not be designed for the specific needs that a NoSql database have. If you want to use a NoSql database as a SQL database, Anvil will probably fit your needs.

### Why does Anvil hates creativity, free thinking and innovation?

We don't hate it, but the main goal of any company is to serve their clients well, and all things that doesn't reach of this goal is a waste of time and money.

Innovation mainly, should not be part of the crucial systems of your company. Have a business is already risk and challenging enough already, you should keep things safe as much as you can, and not try to be different and try to innovate on the authentication service of your system. It's a crucial service, it MUST work 24/7 without bugs.

For a while, companies had the philosophy of "let every team work on it's own way", because it was impossible to make everyone work as a single unit, but it leads to a mountain of technical debit, knowledge lost, malpractices and 1001 reinventions of the wheel.

It forces companies to have multiple developers, that spend most of their time doing things to keep the plates spinning, and don't generate any value for the clients.

We all are here to work, to get a job done, and not to play games and do hobby projects. With Anvil, you can centralize and unify how your teams work. All your services will have the same patterns. Want to change anything? Don't do it one by one, do it all at once, this way you know that all of them are secure and reliable.

### Why is Anvil written in Golang?

Because [Henrique Leite](https://www.linkedin.com/in/henriqueleite42/), the initial main contributor, didn't knew how to effectively write Rust at the time, but we do want to migrate it to Rust someday.

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
