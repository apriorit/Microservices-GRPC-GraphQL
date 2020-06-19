#!/usr/bin/env bash

# Use this script to generate code for micro-services and schemas from for GraphQL form proto-files

# Generate code for services. Code will be placed in gen/go/proto
prototool generate
# Generate schemas. All schemas will be placed in proto/books and proto/holders
protoc --gql_out=paths=source_relative:. -I=. -I=./proto/books ./proto/books/*.proto
protoc --gql_out=paths=source_relative:. -I=. -I=./proto/holders ./proto/holders/*.proto