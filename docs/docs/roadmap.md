---
sidebar_position: 4
---

# Roadmap

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

### CLI Warnings

- Add warning like "Type Foo is unused, maybe you should remove it."
