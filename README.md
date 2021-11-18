# Vagrant Plugin SDK

This repository is a Go library that enables users to write custom [Vagrant](https://vagrantup.com) plugins.

Plugins in Vagrant are separate binaries which communicate with the Vagrant application; the plugin communicates using
gRPC, and while it is theoretically possible to build a plugin in any language supported by the gRPC framework. We
recommend that the developers leverage the [Vagrant SDK](https://github.com/hashicorp/vagrant-plugin-sdk).

## Generating protos

All Go & Ruby protos are wired into `go-generate`. To generate them you'll need a few binaries on your path:

 * `protoc` - installation instructions on the [gRPC Docs](https://grpc.io/docs/protoc-installation/)
 * `grpc_tools_ruby_protoc` - from the [`grpc-tools` gem](https://rubygems.org/gems/grpc-tools/versions/1.41.1), which bundles that binary prebuilt
 * `stringer` - from the [go tools pkg](https://pkg.go.dev/golang.org/x/tools/cmd/stringer)
 * `mockery` - from the [go library hosted at vektra/mockery](https://github.com/vektra/mockery)

You also need to ensure the output directory is present:
 
```sh
$ mkdir -p ruby-proto
```

Once that's all set up you should be ready to roll:

```sh
$ go generate .
```
