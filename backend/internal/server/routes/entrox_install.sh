#!/usr/bin/env bash
set -euo pipefail

APP=entrox
DEFAULT_INSTALL_BASE_URL="https://entrox.996icu.wiki"
INSTALL_DIR=${ENTROX_INSTALL_DIR:-"$HOME/.entrox/bin"}
INSTALL_BASE_URL=${ENTROX_INSTALL_BASE_URL:-$DEFAULT_INSTALL_BASE_URL}
DOWNLOAD_BASE_URL=${ENTROX_DOWNLOAD_BASE_URL:-"${INSTALL_BASE_URL%/}/downloads/entrox-dev"}
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
  -v, --version <version> Install a specific GitHub release version, for example 1.0.180
      --no-modify-path    Do not update shell profile files

Environment:
  ENTROX_DOWNLOAD_BASE_URL  Download mirror base URL. Defaults to ${DEFAULT_INSTALL_BASE_URL}/downloads/entrox-dev
  ENTROX_INSTALL_DIR        Install directory. Defaults to ~/.entrox/bin

Examples:
  curl -fsSL ${DEFAULT_INSTALL_BASE_URL}/install | bash
  ENTROX_DOWNLOAD_BASE_URL=https://your-oss-cdn.example.com/entrox-dev curl -fsSL ${DEFAULT_INSTALL_BASE_URL}/install | bash
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

extract_zip() {
  local archive=$1
  local output=$2

  if command -v unzip >/dev/null 2>&1; then
    unzip -q "$archive" -d "$output"
    return
  fi

  if command -v python3 >/dev/null 2>&1; then
    python3 -m zipfile -e "$archive" "$output"
    return
  fi

  red "Error: 'unzip' or 'python3' is required to extract Entrox."
  red "Ubuntu/Debian: sudo apt-get update && sudo apt-get install -y unzip"
  red "CentOS/RHEL:   sudo yum install -y unzip"
  red "Alpine:       apk add --no-cache unzip"
  exit 1
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

if [[ "$os" == "darwin" && "$arch" == "x64" ]]; then
  rosetta_flag=$(sysctl -n sysctl.proc_translated 2>/dev/null || echo 0)
  if [[ "$rosetta_flag" == "1" ]]; then
    arch=arm64
  fi
fi

case "$os-$arch" in
  darwin-arm64) filename="entrox-cli-macos-arm64.zip" ;;
  linux-x64) filename="entrox-cli-linux-x64.zip" ;;
  windows-x64) filename="entrox-cli-windows-x64.zip" ;;
  darwin-x64)
    red "Intel macOS builds are not available yet. Use Apple Silicon macOS, Linux x64, or Windows x64."
    exit 1
    ;;
  linux-arm64)
    red "Linux arm64 builds are not available yet. Use Linux x64."
    exit 1
    ;;
  *)
    red "Unsupported OS/Arch: $os/$arch"
    exit 1
    ;;
esac

if [[ -n "$VERSION" ]]; then
  download_base="https://github.com/hexonal/entrox/releases/download/v${VERSION#v}"
  display_version="v${VERSION#v}"
else
  download_base="${DOWNLOAD_BASE_URL%/}"
  display_version="latest dev build"
fi

url="$download_base/$filename"
tmp_dir="${TMPDIR:-/tmp}/${APP}_install_$$"
trap 'rm -rf "$tmp_dir"' EXIT
mkdir -p "$tmp_dir" "$INSTALL_DIR"

info "Installing Entrox $display_version for $os-$arch"
info "Downloading $filename"
curl -fL --progress-bar -o "$tmp_dir/$filename" "$url"

extract_zip "$tmp_dir/$filename" "$tmp_dir"

binary=""
for candidate in "$tmp_dir/bin/$APP" "$tmp_dir/$APP" "$tmp_dir/bin/$APP.exe" "$tmp_dir/$APP.exe"; do
  if [[ -f "$candidate" ]]; then
    binary="$candidate"
    break
  fi
done

if [[ -z "$binary" ]]; then
  red "Downloaded archive did not contain the Entrox binary."
  exit 1
fi

install_name="$APP"
if [[ "$os" == "windows" ]]; then
  install_name="$APP.exe"
fi

mv -f "$binary" "$INSTALL_DIR/$install_name"
chmod 755 "$INSTALL_DIR/$install_name"

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
info "Entrox installed to $INSTALL_DIR/$install_name"
info "Run: entrox login"
