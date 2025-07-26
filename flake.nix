{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            caddy
            tailwindcss
            ffmpeg
          ];
        };

        apps.default = {
          type = "app";
          program = toString (pkgs.writeShellScript "serve" ''
            set -x
            rm -r public/
            cp -r static/ public/
            ${pkgs.go}/bin/go run content.go
            ${pkgs.tailwindcss}/bin/tailwindcss -i input.css -o public/css/tailwind.css
            ${pkgs.caddy}/bin/caddy run --config Caddyfile
          '');
        };
      }
    );
}