# rpcguard

`rpcguard` is collection of connect-RPC usage linters which check if connect-RPC method is implemented properly.

- rpc_callvalidate: check if RPC method uses Validate method properly

## Config

- `rpc_callvalidate` provides options. Please see [callvalidate/config.go](passes/callvalidate/config.go)

You can overwrite via commandline option or golangci setting.

## Install

```shell
$ go install github.com/cloverrose/rpcguard/cmd/callvalidate@latest
```

### Or Build from source

```shell
$ make build
```

### Or Install via aqua

https://aquaproj.github.io/

## Usage

```shell
$ go vet -vettool=`which rpc_callvalidate` ./...
```

When you specify config

```shell
go vet -vettool=`which rpc_callvalidate` \
  -rpc_callvalidate.LogLevel=ERROR \
   ./...
```

### Or golangci-lint custom plugin

https://golangci-lint.run/plugins/module-plugins/

Here are reference settings

`.custom-gcl.yml`

```yaml
version: v1.63.4
plugins:
  - module: 'github.com/cloverrose/rpcguard'
    import: 'github.com/cloverrose/rpcguard'
    version: v0.4.0
```

`.golangci.yml`

```yaml
linters-settings:
  custom:
    rpc_callvalidate:
      type: "module"
      description: check if RPC method uses Validate method properly.
      settings:
        LogLevel: "ERROR"
```
