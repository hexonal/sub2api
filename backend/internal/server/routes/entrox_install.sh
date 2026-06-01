#!/usr/bin/env bash
set -euo pipefail

APP=entrox
REPO=hexonal/entrox
INSTALL_DIR=${ENTROX_INSTALL_DIR:-"$HOME/.entrox/bin"}
VERSION=${VERSION:-}
NO_MODIFY_PATH=${NO_MODIFY_PATH:-}

red() {
  printf '\033[0;31m%s\033[0m\n' "$*" >&2
}

info() {
  printf '%s\n' "$*"
}

usage() {
  cat <<EOF
Entrox Installer

Usage: install [options]

Options:
  -h, --help              Display this help message
  -v, --version <version> Install a specific version, for example 1.0.180
      --no-modify-path    Do not update shell profile files

Examples:
  curl -fsSL https://entrox.996icu.wiki/install | bash
  curl -fsSL https://entrox.996icu.wiki/install | bash -s -- --version 1.0.180
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    -h|--help)
      usage
      exit 0
      ;;
    -v|--version)
      if [[ -z "${2:-}" ]]; then
        red "Error: --version requires a value"
        exit 1
      fi
      VERSION="${2#v}"
      shift 2
      ;;
    --no-modify-path)
      NO_MODIFY_PATH=1
      shift
      ;;
    *)
      red "Warning: ignoring unknown option '$1'"
      shift
      ;;
  esac
done

require_command() {
  if ! command -v "$1" >/dev/null 2>&1; then
    red "Error: '$1' is required but not installed."
    exit 1
  fi
}

require_command curl

raw_os=$(uname -s)
case "$raw_os" in
  Darwin*) os=darwin ;;
  Linux*) os=linux ;;
  MINGW*|MSYS*|CYGWIN*) os=windows ;;
  *)
    red "Unsupported OS: $raw_os"
    exit 1
    ;;
esac

arch=$(uname -m)
case "$arch" in
  x86_64|amd64) arch=x64 ;;
  arm64|aarch64) arch=arm64 ;;
  *)
    red "Unsupported architecture: $arch"
    exit 1
    ;;
esac

if [[ "$os" == "windows" && "$arch" != "x64" ]]; then
  red "Unsupported OS/Arch: $os/$arch"
  exit 1
fi

if [[ "$os" == "darwin" && "$arch" == "x64" ]]; then
  rosetta_flag=$(sysctl -n sysctl.proc_translated 2>/dev/null || echo 0)
  if [[ "$rosetta_flag" == "1" ]]; then
    arch=arm64
  fi
fi

needs_baseline=false
if [[ "$arch" == "x64" ]]; then
  if [[ "$os" == "linux" ]]; then
    if ! grep -qwi avx2 /proc/cpuinfo 2>/dev/null; then
      needs_baseline=true
    fi
  elif [[ "$os" == "darwin" ]]; then
    avx2=$(sysctl -n hw.optional.avx2_0 2>/dev/null || echo 0)
    if [[ "$avx2" != "1" ]]; then
      needs_baseline=true
    fi
  elif [[ "$os" == "windows" ]]; then
    ps='(Add-Type -MemberDefinition "[DllImport(""kernel32.dll"")] public static extern bool IsProcessorFeaturePresent(int ProcessorFeature);" -Name Kernel32 -Namespace Win32 -PassThru)::IsProcessorFeaturePresent(40)'
    out=""
    if command -v powershell.exe >/dev/null 2>&1; then
      out=$(powershell.exe -NoProfile -NonInteractive -Command "$ps" 2>/dev/null || true)
    elif command -v pwsh >/dev/null 2>&1; then
      out=$(pwsh -NoProfile -NonInteractive -Command "$ps" 2>/dev/null || true)
    fi
    out=$(echo "$out" | tr -d '\r' | tr '[:upper:]' '[:lower:]' | tr -d '[:space:]')
    if [[ "$out" != "true" && "$out" != "1" ]]; then
      needs_baseline=true
    fi
  fi
fi

is_musl=false
if [[ "$os" == "linux" ]]; then
  require_command tar
  if [[ -f /etc/alpine-release ]]; then
    is_musl=true
  elif command -v ldd >/dev/null 2>&1 && ldd --version 2>&1 | grep -qi musl; then
    is_musl=true
  fi
else
  require_command unzip
fi

target="$os-$arch"
if [[ "$needs_baseline" == "true" ]]; then
  target="$target-baseline"
fi
if [[ "$is_musl" == "true" ]]; then
  target="$target-musl"
fi

archive_ext=.zip
if [[ "$os" == "linux" ]]; then
  archive_ext=.tar.gz
fi
filename="$APP-$target$archive_ext"

if [[ -z "$VERSION" ]]; then
  version=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | sed -n 's/.*"tag_name": *"v\([^"]*\)".*/\1/p' | head -n 1)
  if [[ -z "$version" ]]; then
    red "Failed to resolve latest Entrox release version."
    exit 1
  fi
else
  version="$VERSION"
fi

url="https://github.com/$REPO/releases/download/v$version/$filename"
tmp_dir="${TMPDIR:-/tmp}/${APP}_install_$$"
trap 'rm -rf "$tmp_dir"' EXIT
mkdir -p "$tmp_dir" "$INSTALL_DIR"

info "Installing Entrox $version for $target"
curl -fL --progress-bar -o "$tmp_dir/$filename" "$url"

if [[ "$os" == "linux" ]]; then
  tar -xzf "$tmp_dir/$filename" -C "$tmp_dir"
else
  unzip -q "$tmp_dir/$filename" -d "$tmp_dir"
fi

if [[ ! -f "$tmp_dir/$APP" ]]; then
  red "Downloaded archive did not contain '$APP'."
  exit 1
fi

mv "$tmp_dir/$APP" "$INSTALL_DIR/$APP"
chmod 755 "$INSTALL_DIR/$APP"

add_to_path() {
  local file=$1
  local line=$2
  if [[ -f "$file" && -w "$file" ]] && ! grep -Fxq "$line" "$file"; then
    printf '\n# entrox\n%s\n' "$line" >> "$file"
    info "Added Entrox to PATH in $file"
    return 0
  fi
  return 1
}

if [[ -z "$NO_MODIFY_PATH" ]]; then
  path_line="export PATH=$INSTALL_DIR:\$PATH"
  shell_name=$(basename "${SHELL:-sh}")
  case "$shell_name" in
    zsh) candidates=("${ZDOTDIR:-$HOME}/.zshrc" "${ZDOTDIR:-$HOME}/.zshenv") ;;
    bash) candidates=("$HOME/.bashrc" "$HOME/.bash_profile" "$HOME/.profile") ;;
    fish) candidates=("$HOME/.config/fish/config.fish"); path_line="fish_add_path $INSTALL_DIR" ;;
    *) candidates=("$HOME/.profile" "$HOME/.bashrc") ;;
  esac
  for file in "${candidates[@]}"; do
    if add_to_path "$file" "$path_line"; then
      break
    fi
  done
fi

info ""
info "Entrox installed to $INSTALL_DIR/$APP"
info "Run: entrox login"
