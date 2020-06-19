# How to install utils for code generation

Code generation based on proto-files. All proto-files are processed by [prototool](https://github.com/uber/prototool) that depends on protoc, protoc-gen-go and protoc-gen-gql.
Your machine requires the following utilities to be installed:

1. Install prototool

``` bash
curl -sSL \
  https://github.com/uber/prototool/releases/download/v1.8.0/prototool-$(uname -s)-$(uname -m) \
  -o /usr/local/bin/prototool && \
  chmod +x /usr/local/bin/prototool
```

2. Install protoc

Download the latest relsease from [here](https://github.com/protocolbuffers/protobuf/releases).

Extract archive in a convenient place, go to root of an extracted folder and execute successively:

``` bash
./configure
make
make check
sudo make install
sudo ldconfig
```

3. Install protoc-gen-go

``` bash
go install google.golang.org/protobuf/cmd/protoc-gen-go
```

4. Install protoc-gen-gql

``` bash
go install github.com/danielvladco/go-proto-gql/protoc-gen-gql
```

To generate code for micro-services and also schemas from proto-files use [generate.sh](generate.sh)