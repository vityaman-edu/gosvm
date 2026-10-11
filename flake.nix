{
  description = "Development environment for GoSVM";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs =
    { nixpkgs, ... }:
    let
      supportedSystems = [
        "aarch64-darwin"

        "aarch64-linux"
        "x86_64-linux"
      ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
    in
    {
      devShells = forAllSystems (
        system:
        let
          pkgs = import nixpkgs { inherit system; };
        in
        {
          default = pkgs.mkShell {
            packages = with pkgs; [
              actionlint
              bash
              deadnix
              go
              gofumpt
              golangci-lint
              gopls
              gotestsum
              gotools
              govulncheck
              markdownlint-cli2
              nil
              nixd
              nixfmt
              pre-commit
              shellcheck
              shfmt
              statix
              yamllint
            ];
          };
        }
      );

      formatter = forAllSystems (system: (import nixpkgs { inherit system; }).nixfmt);
    };
}
