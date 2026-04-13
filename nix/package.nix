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

  vendorHash = "sha256-8bHyUrxqNDt6kFwavQLcX3kKOIekea6qQAWSD1KXhfo=";

  proxyVendor = true;

  env.CGO_ENABLED = 0;
  env.GOAMD64 = "v3";

  ldflags = [
    "-s"
    "-w"
    "-X main.version=0.1.0"
  ];

  meta = {
    description = "TUI session manager for Zellij";
    homepage = "https://github.com/pageton/zellij-tui";
    license = lib.licenses.mit;
    mainProgram = "zellij-tui";
  };
}
