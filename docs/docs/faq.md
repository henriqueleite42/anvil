---
sidebar_position: 2
---

# F.A.Q.

## Is Anvil ready to use right out the box?

If you are whiling to adopt one of the architectures proposed by one of ht existing generators, including the ones that we maintain, **yes!**

If you want to implement your own architecture, with your own code patterns, so you will have to write your own generator, what is veeeery complex and time consuming compared to write a single micro-service, but pays off on the long run if you plan to have more than a dozen of micro-services.

## Is Anvil a silver bullet for every project?

No. We are very clear that Anvil is focused on **Event-Oriented, Domain-Driven, Decoupled MicroServices, with a Delivery-Usecase-Repository architecture and SQL Databases**. It may be used in other cases, but we don't give support for these other cases.

## Can Anvil be used for monoliths?

Sure, it works great for both monoliths and micro-services, you only need to find the right generator.

## Can Anvil be used with NoSql databases?

You kinda can if you have the right plugin, but the schema is not and will not be designed for the specific needs that a NoSql database have. If you want to use a NoSql database as a SQL database, Anvil will probably fit your needs.

## Can Anvil be used in a large organization?

Yes, large organizations are the ones that most benefit from Anvil.

In large organizations, we usually have hundreds or even thousands of domains, micro-services, teams, events, packages and changes in the code in a daily basis. It's very hard and demanding to maintain everything, to share these knowledge of the best practices, to ensure that all developers not only know how to implement certain patterns but know the way that the company implements certain patterns.

Anvil works perfectly for these kind of ecosystems. It allows you to have one centralized small team of extremely capable developers that say how the things will work, define rules, best practices, standard libraries, and everything else that you need or want to defined, and all the other teams and members of your organization will follow these rules and patterns.

## Why does Anvil hates creativity, free thinking and innovation?

We don't hate it, but the main goal of any company is to serve their clients well, and all things that doesn't reach of this goal is a waste of time and money.

Innovation mainly, should not be part of the crucial systems of your company. Have a business is already risk and challenging enough already, you should keep things safe as much as you can, and not try to be different and try to innovate on the authentication service of your system. It's a crucial service, it MUST work 24/7 without bugs.

For a while, companies had the philosophy of "let every team work on it's own way", because it was impossible to make everyone work as a single unit, but it leads to a mountain of technical debit, knowledge lost, malpractices and 1001 reinventions of the wheel.

It forces companies to have multiple developers, that spend most of their time doing things to keep the plates spinning, and don't generate any value for the clients.

We all are here to work, to get a job done, and not to play games and do hobby projects. With Anvil, you can centralize and unify how your teams work. All your services will have the same patterns. Want to change anything? Don't do it one by one, do it all at once, this way you know that all of them are secure and reliable.
