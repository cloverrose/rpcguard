# rpcguard

`rpcguard` is collection of connect-RPC usage linters which check if connect-RPC method is implemented properly.

- rpc_callvalidate: check if RPC method uses Validate method properly
- rpc_wraperr: check if RPC method returns wrapped error

## Config

- `rpc_callvalidate` provides options. Please see [callvalidate/config.go](passes/callvalidate/config.go)
- `rpc_wraperr` provides options. Please see [wraperr/config.go](passes/wraperr/config.go)

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

```shell
$ go vet -vettool=`which rpc_wraperr` -rpc_wraperr.IncludePackages="$(go list -m)/.*" ./...
```

Note: rpc_wraperr.IncludePackages is required option.


When you specify config

```shell
go vet -vettool=`which rpc_callvalidate` \
  -rpc_callvalidate.LogLevel=ERROR \
   ./...
```

```shell
go vet -vettool=`which rpc_wraperr` \
  -rpc_wraperr.IncludePackages="$(go list -m)/.*" \
  -rpc_wraperr.LogLevel=ERROR \
  -rpc_wraperr.ReportMode=RETURN \
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
    rpc_wraperr:
      type: "module"
      description:  check if RPC method returns wrapped error.
      settings:
        LogLevel: "ERROR"
        ReportMode: "RETURN"
        IncludePackages: "github.com/cloverrose/linterplayground/.*"
```
