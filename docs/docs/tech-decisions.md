---
sidebar_position: 49
---

# Technical Decisions

See our technical decisions and why we took them.

### Why is Anvil written in Golang?

Because [Henrique Leite](https://henriqueleite42.com), the initial main contributor, didn't knew how to effectively write Rust at the time, but we do want to migrate it to Rust someday.

### Why YAML?

The options that we have are: YAML, JSON, TOML and INI. YAML is the best between the three for long nested configuration. TOML doesn't look good and can be extremely repetitive when dealing with these things, JSON doesn't support comments and INI is more focused in configuration, not specification.

### Why does Anvil have some hidden behavior?

Anvil sometimes can have hidden behaviors, for example, if you have a HTTP delivery that returns an StatusCode 204, the body will not be
returned, and this is not explicit in anywhere in the config files. We do it because our focus here is to help developers to have best practices
while having to write the least amount of code.

In HTTP status code specification, it's very explicit that StatusCode 204 should not return any content, so instead of having a way to configure
it, Anvil already handles it for you, ensuring that you follow best practices.

### Problems with ordering

To have the best UX writing the schema in Yaml, we opted to use Maps instead of Lists in many of the cases, here's an example:

With Maps, the current structure of `Entities` is this:
```yaml
Entities:
  Entities:
    User:
      DbName: users
      Columns:
        Id:
          Type: Int
```

And with Lists, it would be like this:
```yaml
Entities:
  Entities:
    - Name: User
      DbName: users
      Columns:
        - Name: Id
          Type: Int
```

It doesn't change that much, but we fear that the `$ref` would be confuse, because you must use the `Name` and not the `index`, what may be confusing.

Using Maps has a problem: We are unable to keep the order of the things, due to limitations of the data structure. We could create our own parse for the YAML file, but the YAML specification is too complex for us to do it while maintaining Anvil and it's generators, so we opted for an work around, the `Order` prop:

```yaml
Entities:
  Entities:
    User:
      Name: users
			Order: 0
      Columns:
        Id:
          Type: Int
					Order: 0
```

If you don't want to use the default ordering (alphabetical), you can set the order manually. It doesn't need to follow a specific order, it can have any value that you want, as long as it's an `int32`.

If you are using grpc and has deprecated fields, you can use `Order` to ignore these fields, like this:

```yaml
Entities:
  Entities:
    User:
      Name: users
			Order: 0
      Columns:
        Id:
          Type: Int
					Order: 0
        Foo:
          Type: String
					Order: 1
        Bar:
          Type: String
					Order: 3 # Not a problem, will work fine
```
