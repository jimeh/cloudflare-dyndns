#! /usr/bin/env bash
set -e
shopt -s extglob

main() {
  local name="cloudflare-dyndns"
  local platforms=(
    darwin-386 darwin-amd64
    freebsd-386 freebsd-amd64 freebsd-arm
    linux-386 linux-amd64 linux-arm
    netbsd-386 netbsd-amd64 netbsd-arm
    openbsd-386 openbsd-amd64
    solaris-amd64
    windows-386 windows-amd64
  )
  local builddir="build"
  local outputdir="pkg"

  local version
  version="$(get-version)"

  local workdir="${builddir}/${version}"
  mkdir -p "$workdir"

  for platform in "${platforms[@]}"; do
    if [[ "$platform" =~ ^(.+)-(.+)$ ]]; then
      local os="${BASH_REMATCH[1]}"
      local arch="${BASH_REMATCH[2]}"
      local pkg="${name}_${version}_${os}_${arch}"
      local pkgdir="${workdir}/${pkg}"
      local binary="${pkgdir}/${name}"

      if [ "$os" == "windows" ]; then
        binary="${binary}.exe"
      fi

      echo "building $pkg"
      GOOS="$os" GOARCH="$arch" go build -o "$binary"

      cp "README.md" "${pkgdir}/"

      mkdir -p "${outputdir}/${version}"

      if [ "$os" == "windows" ]; then
        local archive="${outputdir}/${version}/${pkg}.zip"
        echo "creating ${archive}"
        local cwd="$(pwd)"
        cd "$workdir"
        zip -r "../../$archive" "$pkg"
        cd "$cwd"
      else
        local archive="${outputdir}/${version}/${pkg}.tar.gz"
        echo "creating ${archive}"
        tar -C "$workdir" -czf "$archive" "$pkg"
      fi
    fi
  done
}

get-version() {
  trim "$(cat "VERSION")"
}

trim() {
  local string="$@"
  string="${string#"${string%%[![:space:]]*}"}"
  string="${string%"${string##*[![:space:]]}"}"
  echo -n "$string"
}

main
