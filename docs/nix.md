# Nix Setup

zellij-tui ships a Nix flake with pre-configured outputs for building, development, and declarative Home Manager configuration.

## Prerequisites

| Requirement | Version  | Notes                                    |
|-------------|----------|------------------------------------------|
| Nix         | 2.4+     | With flakes enabled                      |
| Flakes      | enabled  | See below if not already enabled         |

### Enable flakes (if needed)

Add to `/etc/nix/nix.conf` or `~/.config/nix/nix.conf`:

```
experimental-features = nix-command flakes
```

On NixOS, set in your configuration:

```nix
nix.settings.experimental-features = [ "nix-command" "flakes" ];
```

## Quick start

### Build the package

```bash
nix build github:pageton/zellij-tui
./result/bin/zellij-tui
```

### Run directly (no install)

```bash
nix run github:pageton/zellij-tui
```

### Enter the dev shell

```bash
nix develop github:pageton/zellij-tui
```

Provides: `go`, `gopls`, `gotools`, `zellij`.

## Flake outputs

| Output                          | Description                                  |
|---------------------------------|----------------------------------------------|
| `packages.<system>.default`     | The `zellij-tui` binary                      |
| `packages.<system>.zellij-tui`  | Same as `default` (explicit name)            |
| `devShells.<system>.default`    | Dev shell with Go toolchain + zellij         |
| `overlays.default`              | Overlay adding `zellij-tui` to `pkgs`        |
| `homeManagerModules.default`    | Home Manager module (`programs.zellij-tui`)  |

Supported systems: `x86_64-linux`, `aarch64-linux`, `x86_64-darwin`, `aarch64-darwin`.

---

## Home Manager module

The flake provides a Home Manager module that installs zellij-tui and generates its config file declaratively from Nix.

### Basic setup

In your Home Manager config (`home.nix`, `home-manager/home.nix`, or a module):

```nix
{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    home-manager = {
      url = "github:nix-community/home-manager";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    zellij-tui = {
      url = "github:pageton/zellij-tui";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { nixpkgs, home-manager, zellij-tui, ... }: {
    homeConfigurations.your-user = home-manager.lib.homeManagerConfiguration {
      pkgs = nixpkgs.legacyPackages.x86_64-linux;
      modules = [
        ./home.nix
        zellij-tui.homeManagerModules.default
      ];
    };
  };
}
```

Then in your `home.nix`:

```nix
{ pkgs, ... }:

{
  programs.zellij-tui = {
    enable = true;

    settings = {
      zellij_path = "${pkgs.zellij}/bin/zellij";
      auto_attach = true;
    };
  };
}
```

This does two things:

1. Adds `zellij-tui` to your `PATH`
2. Generates `~/.config/zellij-tui/config.toml` from the `settings` attrset

### Module options

#### `programs.zellij-tui.enable`

- **Type:** `boolean`
- **Default:** `false`

Enable zellij-tui. When enabled, the package is added to `home.packages`.

#### `programs.zellij-tui.package`

- **Type:** `package`
- **Default:** the package built from this flake

Override to use a custom/patched version:

```nix
programs.zellij-tui = {
  enable = true;
  package = pkgs.callPackage ./my-custom-zellij-tui.nix { };
};
```

#### `programs.zellij-tui.settings`

- **Type:** attribute set (TOML-compatible)
- **Default:** `{ }`

Settings written to `$XDG_CONFIG_HOME/zellij-tui/config.toml`. Null values are filtered out.

##### `settings.zellij_path`

- **Type:** `nullOr str`
- **Default:** `null` (PATH lookup at runtime)

Path to the zellij binary. When `null`, zellij-tui finds zellij via `PATH`. Set this to pin to a specific nixpkgs zellij:

```nix
zellij_path = "${pkgs.zellij}/bin/zellij";
```

##### `settings.auto_attach`

- **Type:** `bool`
- **Default:** `true`

Immediately attach to a newly created session. When `false`, sessions are created in the background and the TUI stays open.

### Full example

```nix
{ pkgs, ... }:

{
  programs.zellij = {
    enable = true;
  };

  programs.zellij-tui = {
    enable = true;

    settings = {
      zellij_path = "${pkgs.zellij}/bin/zellij";
      auto_attach = true;
    };
  };

  # Optional: auto-launch on terminal open
  programs.bash.initExtra = ''
    if [[ -z "$ZELLIJ" ]] && [[ $- == *i* ]]; then
      zellij-tui
    fi
  '';

  programs.zsh.initExtra = ''
    if [[ -z "$ZELLIJ" ]] && [[ $- == *i* ]]; then
      zellij-tui
    fi
  '';
}
```

---

## Using the overlay

If you want `zellij-tui` available in your nixpkgs as `pkgs.zellij-tui`:

```nix
nixpkgs.overlays = [
  zellij-tui.overlays.default
];
```

Then use it anywhere:

```nix
environment.systemPackages = [ pkgs.zellij-tui ];
# or
home.packages = [ pkgs.zellij-tui ];
```

---

## direnv integration

If you use [direnv](https://direnv.net), add an `.envrc` to the project root:

```bash
# .envrc
use flake
```

Then:

```bash
echo "use flake" > .envrc
direnv allow
```

This drops you into the dev shell automatically whenever you `cd` into the project directory.

---

## Building from a local clone

If you're hacking on zellij-tui:

```bash
git clone https://github.com/pageton/zellij-tui.git
cd zellij-tui

# Build
nix build

# Dev shell
nix develop

# Run tests
nix develop -c go test ./...

# Run
nix run
```

### Updating the vendor hash

If you change Go dependencies (`go get`, `go mod tidy`), the vendor hash must be updated in `nix/package.nix`:

```bash
# Build will fail and print the correct hash
nix build

# Copy the "got:" hash from the error and update vendorHash in nix/package.nix
```

Or use `nix-prefetch`:

```bash
nix-prefetch -I nixpkgs=flake:nixpkgs '{ sha256 }: (import ./. {}).packages.x86_64-linux.zellij-tui.goModules.overrideAttrs (_: { vendorHash = sha256; })'
```

---

## NixOS module (standalone, no Home Manager)

If you don't use Home Manager and want a system-level install:

```nix
{ pkgs, ... }:

let
  zellij-tui = builtins.getFlake "github:pageton/zellij-tui";
in
{
  environment.systemPackages = [
    zellij-tui.packages.${pkgs.stdenv.hostPlatform.system}.default
    pkgs.zellij
  ];
}
```

Config would need to be managed manually (create `~/.config/zellij-tui/config.toml` yourself), or via `environment.etc` / `systemd.tmpfiles`.

---

## Troubleshooting

### `error: hash mismatch in fixed-output derivation`

The `vendorHash` in `nix/package.nix` is stale. Rebuild to get the correct hash from the error output, then update the file.

### `warning: Git tree is dirty`

You have uncommitted changes. Either commit them or use `--impure`:

```bash
nix build --impure
```

### Flakes not finding the repo

Make sure you're running from inside the git repo and all files are tracked (`git add`):

```bash
git add -A
nix build
```

### `experimental-features 'flakes' is not enabled`

See [Enable flakes](#enable-flakes-if-needed) above.
