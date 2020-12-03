# Vagrant Plugin SDK

This repository is a Go library that enables users to write custom [Vagrant](https://vagrantup.com) plugins.

Plugins in Vagrant are separate binaries which communicate with the Vagrant application; the plugin communicates using
gRPC, and while it is theoretically possible to build a plugin in any language supported by the gRPC framework. We
recommend that the developers leverage the [Vagrant SDK](https://github.com/hashicorp/vagrant-plugin-sdk).

## Generating protos

To generate go protos run

```
$ go generate .
```

To generate ruby protos run
(Remember to install gems `grpc-tools` and `grpc`)

```
$ sh -c "grpc_tools_ruby_protoc -I `go list -m -f \"{{.Dir}}\" github.com/mitchellh/protostructure` -I ./3rdparty/proto/api-common-protos -I proto --ruby_out=proto/gen/ --grpc_out=proto/gen/ ./proto/plugin.proto"
```
