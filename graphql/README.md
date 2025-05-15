# GraphQL experiment
* ~~Feels complex when I setup with gqlgen (Go)~~ I didn't read docs https://gqlgen.com/getting-started/ and relied on ChatGPT.
  * Understand that not everytime you can get all answers by repeatedly asking question. It's tiresome and you'll miss out on basics
* Define the Schema
  * The schema had model, query, mutation and subscription 
* Use `gqlgen generate` command to generate the graphql resolvers
* The gqlgen configuration is defined in `gqlgen.yml` file. This file has the info on where to read the graphql files, and where to put the generated files
* There is a default server code that handles all the queries.
* Start the server with `go run ./server/`


# GqlGen specifics
* Set up a `Resolver` struct. This is the root resolver. 
  * This has the Query, Mutation and Subscription resolver by default
* The Resolver struct have `viewerResolve`, `mutationResolver`, and `subscriptionResolver` defined and any other Type resolvers if needed
  * You can generate resolvers and say what models needs resolvers in gqlgen.yml
* Then these resolver have the registered methods that will be called.
* GraphQL -> Resolver.TodoResolver -> Todos query -> If user is present, resolve users -> Return the Todos

# Example GraphQL queries
* Goto `localhost:8080`
```graphql
query {
  person(id: "abc-123") {
    id
    name
    age
    email
  }
}

mutation {
  createPerson(name: "Alice", age: 30, email: "alice@example.com") {
    id
    name
    age
    email
  }
}
```

# Something to read on
* Why this guy is over GraphQL - https://bessey.dev/blog/2024/05/24/why-im-over-graphql/
* Persisted Query is a solution for - caching, rate limiting, query control (but make it REST like?)
  * You define and store the query as hash. You use the hash (or id) in the frontend to get the data

# GraphQL Federation
* GraphQL also grows in complexity cuz it units all the services LOL
* Federation means a unified Organization where internal sub org have their own influence. Likewise GraphQL Federation unifies all the micrographqls (like gateway for microservices)
* It doesn't only do that. It can make all the subgraphs look like one. You can have Account in one subgraph and refer to it in another subgraph
* See this repo for more examples - https://github.com/99designs/gqlgen/blob/master/_examples/federation
