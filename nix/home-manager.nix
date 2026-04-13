{
  self,
  config,
  lib,
  pkgs,
  ...
}:
with lib;
let
  cfg = config.programs.zellij-tui;

  settingsFormat = pkgs.formats.toml { };
in {
  options.programs.zellij-tui = {
    enable = mkEnableOption "zellij-tui — TUI session manager for Zellij";

    package = mkOption {
      type = types.package;
      default = self.packages.${pkgs.stdenv.hostPlatform.system}.default;
      defaultText = literalExpression "self.packages.\${system}.default";
      description = "The zellij-tui package to use.";
    };

    settings = mkOption {
      type = types.submodule {
        freeformType = settingsFormat.type;
        options = {
          zellij_path = mkOption {
            type = types.nullOr types.str;
            default = null;
            example = "\${pkgs.zellij}/bin/zellij";
            description = ''
              Path to the zellij binary.
              When null (default), zellij-tui does a PATH lookup at runtime.
              Set this to pin to a specific zellij from nixpkgs.
            '';
          };

          auto_attach = mkOption {
            type = types.bool;
            default = true;
            description = ''
              Whether to immediately attach after creating a new session.
              When false, sessions are created in the background.
            '';
          };
        };
      };
      default = { };
      description = ''
        Configuration written to
        {file}`$XDG_CONFIG_HOME/zellij-tui/config.toml`.
      '';
      example = literalExpression ''
        {
          zellij_path = "''${pkgs.zellij}/bin/zellij";
          auto_attach = true;
        }
      '';
    };
  };

  config = mkIf cfg.enable {
    home.packages = [ cfg.package ];

    xdg.configFile."zellij-tui/config.toml" = mkIf (cfg.settings != { }) {
      source = settingsFormat.generate "zellij-tui-config.toml" (
        filterAttrs (n: v: v != null) cfg.settings
      );
    };
  };
}
