# How to generate GraphQL code from schemas

Copy generated *.graphql files from /proto/books and /proto/holders folders to /schemas and modify them if needed.

Install gqlgen:

``` bash
go get github.com/99designs/gqlgen
```

Go to folder with schemas:

``` bash
cd graph_api/schemas
```

Run code generation:

``` bash
gqlgen generate
```
