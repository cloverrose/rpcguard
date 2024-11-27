# rpcguard

`rpcguard` checks if connect-RPC endpoint method is implemented properly.

## Config

`rpcguard` provides options. Please see [config.go](config.go)

You can overwrite via commandline option or golangci setting.

- Verbose

## Install

```shell
$ go install github.com/cloverrose/rpcguard/cmd/rpcguard@latest
```

### Or Build from source

```shell
$ go build -o bin/ ./cmd/...
```

### Or Install via aqua

https://aquaproj.github.io/

## Usage

```shell
$ go vet -vettool=`which rpcguard` ./...
```

When you specify config

```shell
go vet -vettool=`which rpcguard` \
  -rpcguard.Verbose=true \
   ./...
```

### Or golangci-lint custom plugin

https://golangci-lint.run/plugins/module-plugins/

Here are reference settings

`.custom-gcl.yml`

```yaml
version: v1.62.0
plugins:
  - module: 'github.com/cloverrose/rpcguard'
    import: 'github.com/cloverrose/rpcguard'
    version: v0.1.0
```

`.golangci.yml`

```yaml
linters-settings:
  custom:
    rpcguard:
      type: "module"
      description: check connect-RPC endpoint implementation.
      settings:
        Verbose: true
```
