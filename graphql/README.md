# GraphQL experiment
* Feels complex when I setup with gqlgen (Go)
* I had to define a schema
  * The schema had model, query, mutation and subscription 
* I used `gqlgen` command to generate the graphql resolvers
* The gqlgen configuration is defined in `gqlgen.yml` file. This file has the info on where to read the graphql files, and where to put the generated files
* There is a default server code that handles all the queries.
* Start the server with `go run ./server/`


# GqlGen specifics
* Set up a `Resolver` struct.
* The Resolver struct have `viewerResolve`, `mutationResolver`, and `subscriptionResolver` defined
* Then these resolver have the registered methods that will be called.
* GraphQL -> Resolver.viewerResolver -> person query -> Return the person

# Example GraphQL queries
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
