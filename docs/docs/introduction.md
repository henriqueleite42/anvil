---
sidebar_position: 1
---

# Introduction

Anvil is a boilerplate generator for APIs. It focuses on unifying all the parts of your API in one single source of truth and allowing you to change architectures, frameworks, and libraries with minimum effort.

It follows a  _schema-first_ approach, instead of writing your code first, you write a schema, and it generates most of the code for you, all the parts that are repetitive and don't influence the performance, while letting you have 100% control of the part that matters: the business logic.

It can be used for both micro-services and monoliths.

## Why use Anvil

Anvil will help you to:
- Need less developers to accomplish the same (probably even best) results
- Decrease the amount of time that it takes to create new products and features, without having to compromise the quality and security of the software
- Better divide the responsibilities of your team, to get the best that they can offer and not needing so many experienced developers to create amazing products
- Document you code, APIs and database with almost no effort

## How Anvil does these things?

Anvil by itself:
- _Schema-first_ approaches help you to visualize the current state of your system in a very easy and fast way: Instead of having to understand code, the project pattern, searching in a bunch of files, go directly to the ONE file definition and figure it out right away.
- Instead of trying to reinvent the wheel, you can follow a standardized architecture that is scalable, clean, flexible, follows the best practices and allows the work to be divided in multiple steps that can be executed in parallel
- Standardize all you domains / services to follow the exact same patterns for EVERYTHING, decreasing a lot the learning curve and the effort necessary to maintain them
- Allows tracking and usage of confidential and private data, like user's emails, to complain with regulations

Anvil generators:
- Generate e2e tests, useful for early stage startups that can't afford a QA or have enough time to implement more complex tests
- Generate `.proto` files for gRPC APIs and OpenAPI specs for REST APIs
- Generate database migrations and automatically handle them, begin able to easily see the current state of your database by looking at the schema file
- Generate standardized clients for your APIs, with automatically generate `CHANGELOG.md`s that follows [SemVer](https://semver.org)
- You are not stuck to Anvil. It's not a framework, it generates code that you have 100% control of. If you don't want to use Anvil anymore, it will have no impact on your systems.
- Easy to change things at a global scale: Do not get stuck into a language, framework or architecture anymore, if you want to change it one day, it will be 100x times easier and faster.
- Run them in the pipeline to automatically update things without the need of a developer

Anvil platform (coming soon):
- Easily visualize your database and documentations
- Automatically notify all the teams dependent on your API about updates
- Integration with other tools like Jira and Linear, to automatically create tasks when the schema changes

## F.A.Q.

[See the F.A.Q. here.](./faq)

## Examples

[See the examples here.](https://github.com/henriqueleite42/anvil/cli/tree/master/examples)

## How it works

[Understand how it works here.](./how-it-works)
