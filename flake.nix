{
  description = "Terminal User Interface for managing tmux'es windows and sessions";

  inputs =
    {
      nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

      flake-utils.url = "github:numtide/flake-utils";

      gitignore.url = "github:hercules-ci/gitignore.nix";
      gitignore.inputs.nixpkgs.follows = "nixpkgs";
    };

  outputs = inputs:
    let
      inherit (inputs) nixpkgs gitignore flake-utils;
      inherit (gitignore.lib) gitignoreSource;
    in
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      rec {
        formatter = pkgs.nixpkgs-fmt;
        packages.default = packages.tmux-tui;
        packages.tmux-tui = pkgs.buildGoModule {
          pname = "tmux-tui";
          version = (builtins.readFile ./tmux_tui/version);
          src = gitignoreSource ./.;
          vendorHash = "sha256-AmBosdrk/pzqJrRAhxhVwLjceKTdh07nTdVonVUTa/A=";
          installPhase = ''
            runHook preInstall
            mkdir -p $out/bin
            mkdir -p build
            $GOPATH/bin/docgen
            cp -r build/share $out/share
            cp $GOPATH/bin/tmux-tui $out/bin/tmux-tui
            strip $out/bin/tmux-tui
            runHook postInstall
          '';
        };
        apps = rec {
          tmux-tui = { type = "app"; program = "${packages.tmux-tui}/bin/tmux-tui"; };
          default = tmux-tui;
        };
        devShell = pkgs.mkShell {
          packages = with pkgs;[ packages.tmux-tui go git man busybox ];
          shellHook = ''
            export fish_complete_path=${packages.tmux-tui}/share/fish/completions
          '';
        };
      }
    );
}
