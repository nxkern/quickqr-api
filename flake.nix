{
  description = "quickqr";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      nixpkgs,
      flake-utils,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };

        packages = [
          pkgs.go
          pkgs.gopls
          pkgs.gotestsum
          pkgs.golangci-lint
          pkgs.air
          pkgs.delve
          pkgs.govulncheck
          pkgs.dive
          pkgs.lefthook
        ];
      in
      {
        devShells.default = pkgs.mkShell {
          inherit packages;
          shellHook = ''
            export CGO_ENABLED=0
            echo -e "\033[1;32m✔ quickqr-api environment active.\033[0m"
          '';

        };

      }
    );
}
