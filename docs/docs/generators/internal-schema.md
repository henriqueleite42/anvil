---
sidebar_position: 2
---

# Internal Schema

We have 2 `schema`s.

The first thing that you need to know when developing for Anvil is that the `schema` that Anvil uses internally is completely different than the one that the users write (`.anv` files).

While the `schema` for written by users is focused on:
- Being easily understandable for humans
- Being able to share knowledge to other humans (trough comments)
- Being easy and lean to write
- Being divisible in multiple files
- Organized in the way that you architect your API

The internal `schema` that Anvil uses is focused on:
- Being performative to parse by the machine
- Having everything that you need in one single file
- Already have all the references resolved
- Most performative as possible to access all the values that you may need
- Organized on the way that your program will generate things

To identify the difference between processed / internal `schema`s and "user-friendly" `schema`s, we have 2 different types of files:
- `.anv`: user friendly schema, the one described [here](../use/schema.md).
- `.anvp`: processed schema, made by the machine for the machine

You will never touch a `.anvp` file, but Anvil will always use this format to communicate to your generator.

## When and how does Anvil CLI communicates with the generator?

Anvil CLI calls your generator when the configuration says so.

Anvil CLI uses `stdout`, a decision inspired by LSPs. All the information is formatted in **JSON string**.
