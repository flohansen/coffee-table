{
  description = "Coffee Table";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem(system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
          go-migrate = pkgs.go-migrate.overrideAttrs(oldAttrs: {
            tags = ["sqlite"];
          });
        in
        {
          devShells.default = pkgs.mkShell {
            buildInputs = with pkgs; [
              go
              go-migrate
              sqlc
            ];
          };
        });
}
