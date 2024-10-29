---
sidebar_position: 3
---

# How it works

Anvil has 5 main parts, each one responsible for a specific complementary role.

## `*.anv` files

The schema definition is a `*.anv` file that describes a domain of your project. Each project can have multiple domains in it, and they can be related or not.

Think about the `.anv` files like a `schemas.prisma` or an OpenApi spec, and from this, we generate a variety of things.

## `anvil.yaml`

`anvil.yaml` is the configuration file for Anvil, where you put the paths for the schemas and configure the generators that you are using.

## CLI

The CLI is the way that you interact with all Anvil things. You can use it to validate your files, generate things, install plugins, run your migrations, and much more.

It's designed to work with CI/CD too 🙌

## Generators

Generators allow you to generate code based on a `*.anv` config. They come in various shapes and sizes, and can be used for practically anything:
- Generate a project with a specific code pattern, that uses a specific set of libraries
- Generate e2e tests
- Generate changelogs

Generators are were the magic of Anvil happens.
