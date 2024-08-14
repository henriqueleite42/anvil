# Our Technical Decisions and how we took them

### Why is Anvil written in Golang?

Because [Henrique Leite](https://henriqueleite42.com), the initial main contributor, didn't knew how to effectively write Rust at the time, but we do want to migrate it to Rust someday.

### Why YAML?

The options that we have are: YAML, JSON and TOML. YAML is the best between the three for long nested configuration. TOML doesn't look good and can be extremely repetitive when dealing with these things, and JSON doesn't support comments.

### Why does Anvil have some hidden behavior?

Anvil sometimes can have hidden behaviors, for example, if you have a HTTP delivery that returns an StatusCode 204, the body will not be
returned, and this is not explicit in anywhere in the config files. We do it because our focus here is to help developers to have best practices
while having to write the least amount of code.

In HTTP status code specification, it's very explicit that StatusCode 204 should not return any content, so instead of having a way to configure
it, Anvil already handles it for you, ensuring that you follow best practices.
