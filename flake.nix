{
  description = "HashiCorp Vagrant SDK";

  inputs.vagrant.url =
    "/home/spox/Projects/hashicorp/vagrant"; # "github:hashicorp/vagrant-ruby";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, flake-utils, vagrant }:
    flake-utils.lib.eachDefaultSystem (system: {
      # Just use the exact same shell environment as Vagrant.
      devShells.default = vagrant.devShells.${system}.default;
    });
}
