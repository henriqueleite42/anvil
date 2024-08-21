---
sidebar_position: 1
---

# Introduction

Thanks for wanting to develop things for Anvil! We are here to try to make your journey as smooth as possible.

In this page, we will guide you about the necessary things that you need to know **BEFORE** starting to create your generator.

## What is a generator?

Generators are CLIs that generate things.

These things can be code, schemas, migrations, anything you want.

They should not be used to perform side-actions, like create tasks on Jira, or do things that aren't related to generate things: Your generator MUST have an output.

To do side actions use [plugins](../plugins/introduction.md).

## Do I need to use a specific language to create?

No, you can use any programming language that you want, you just need to make it executable by the final user.

Example: If you are creating a generator for NodeJs projects using NodeJs, the final user probably will have NodeJs installed, so you can ship JS files directly.

You can run your generator any way that you want, the command to run it is specified on the `.anvilconfig` file.

## How to create a generator?

:::warning

**BEFORE** creating your generator, please get to know and understand our [internal schema](./internal-schema).

:::
