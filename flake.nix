{
  description = "EPUB renamer";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages = {
          default = pkgs.buildGoModule {
            name = "epmv";
            src = self;
            vendorHash = "sha256-fOqtjEoZ6KSryGex/y2mIWejrTGok3GSK2phssRM5qg=";
            goSum = ./go.sum;
          };
        };

        devShells = pkgs.mkShell {
          packages = with pkgs; [
            go
          ];
        };

        formatter = pkgs.nixfmt-tree;
      }
    );
}