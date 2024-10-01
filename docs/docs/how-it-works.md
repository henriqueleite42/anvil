---
sidebar_position: 3
---

# How it works

Anvil has 5 main parts, each one responsible for a specific complementary role.

## `*.anv` files

The schema definition is a `.anv` file that describes a domain of your service. Each project (micro-service) can have multiple domains in it, and they can be related or not (ideally, if they are in the same project, they should be).

Think about the `.anv` files like a `schemas.prisma` or an OpenApi spec, and from this we generate an infinity of things.

## `.anvilconfig`

`.anvilconfig` is the configuration file for Anvil, where you put information like the plugins that you are using, the things that you want to generate, and any other configuration that Anvil CLI or the plugins may need.

It is written in [INI](https://en.wikipedia.org/wiki/INI_file).

## CLI

The CLI is the way that you interact with all Anvil things. You can use it to validate your files, generate things, install plugins, run your migrations, and much more.

It's designed to work with CI/CD too 🙌

## Generators

Generators allows you to generate code based on a `.anv` config. They come in various shapes and sizes, and can be used for practically anything:
- Generate an microservice with a specific code pattern, that uses a specific set of libraries
- Generate e2e tests
- Generate changelogs

Generator are were the magic oh Anvil happens.

## Plugins

Plugins do side-effects with Anvil, like:
- Create tasks in Jira based on the changes of the schemas.
- Notify breaking changes to dependent projects

Very useful for Agile environments.
