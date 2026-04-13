{
  lib,
  buildGoModule,
}:
buildGoModule {
  pname = "zellij-tui";
  version = "0.1.0";

  src = lib.fileset.toSource {
    root = ../.;
    fileset = lib.fileset.unions [
      ../go.mod
      ../go.sum
      ../main.go
      ../internal
    ];
  };

  vendorHash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";

  ldflags = [
    "-s"
    "-w"
  ];

  meta = {
    description = "TUI session manager for Zellij";
    homepage = "https://github.com/sadiq/zellij-tui";
    license = lib.licenses.mit;
    mainProgram = "zellij-tui";
  };
}
