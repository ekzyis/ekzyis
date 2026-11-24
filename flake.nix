{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        build = pkgs.writeShellScript "build" ''
          set -xe
          cp -r static/. public/
          ${pkgs.go}/bin/go run *.go
          ${pkgs.tailwindcss}/bin/tailwindcss -i input.css -o public/css/tailwind.css
        '';
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            caddy
            tailwindcss
            ffmpeg
            s3cmd
            certbot
          ];
        };

        apps.default = {
          type = "app";
          program = toString (pkgs.writeShellScript "serve" ''
            set -xe
            rm -rf public/
            ${build}

            ${pkgs.caddy}/bin/caddy run --config Caddyfile &
            trap "kill $!" EXIT

            ${pkgs.watchexec}/bin/watchexec --postpone --on-busy-update=queue \
              -w content -w html -w static -w input.css -w tailwind.config.js -w xbuild.go -w chroma.go \
              -- ${build}
          '');
        };

        apps.deploy = {
          type = "app";
          program = toString (pkgs.writeShellScript "deploy" ''
            set -e

            s3cmd_opts="--guess-mime-type --no-mime-magic --acl-public --delete-removed --delete-after"

            log () {
                echo "+" "$@"
                "$@"
            }

            log rm -rf public/
            log cp -r static/ public/
            log ${pkgs.go}/bin/go run *.go -M
            log ${pkgs.tailwindcss}/bin/tailwindcss -i input.css -o public/css/tailwind.css

            log ${pkgs.s3cmd}/bin/s3cmd sync $s3cmd_opts --dry-run public/ s3://www.ekzy.is

            echo
            read -n 1 -p "Continue deploy? [yn] " yn
            echo

            [ "$yn" == "y" ] && log ${pkgs.s3cmd}/bin/s3cmd sync $s3cmd_opts public/ s3://www.ekzy.is
          '');
        };
      }
    );
}
