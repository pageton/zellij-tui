{
  description = "zellij-tui — TUI session manager for Zellij";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    home-manager = {
      url = "github:nix-community/home-manager";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      home-manager,
    }:
    let
      forEachSystem = nixpkgs.lib.genAttrs [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
    in
    {
      packages = forEachSystem (system: rec {
        default = zellij-tui;
        zellij-tui = nixpkgs.legacyPackages.${system}.callPackage ./nix/package.nix { };
      });

      devShells = forEachSystem (system: {
        default = nixpkgs.legacyPackages.${system}.mkShell {
          packages = with nixpkgs.legacyPackages.${system}; [
            go
            gopls
            gotools
            zellij
            golangci-lint
            staticcheck
            goimports-tools
            just
          ];

          shellHook = ''
            echo "zellij-tui dev shell — go $(go version)"
          '';
        };
      });

      overlays.default = final: _prev: {
        zellij-tui = self.packages.${final.stdenv.hostPlatform.system}.default;
      };

      homeManagerModules = rec {
        default = zellij-tui;
        zellij-tui = import ./nix/home-manager.nix self;
      };
    };
}
