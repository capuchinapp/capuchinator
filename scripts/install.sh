#!/bin/bash

# This runs on Unix shells like bash/dash/ksh/zsh. It uses the common `local`
# extension. Note: Most shells limit `local` to 1 var per line, contra bash.

# Some versions of ksh have no `local` keyword. Alias it to `typeset`, but
# beware this makes variables global with f()-style function syntax in ksh93.
# mksh has this alias by default.
has_local() {
  # shellcheck disable=SC2034  # deliberately unused
  local _has_local
}
has_local 2>/dev/null || alias local=typeset

need_cmd mkdir
need_cmd curl
need_cmd tar
need_cmd chmod
need_cmd rm

HOME_BIN_DIR="~/.local/bin"
VERSION=v0.6.0
URL="https://github.com/capuchinapp/capuchinator/releases/download/${VERSION}/capuchinator_${VERSION#v}_linux_amd64.tar.gz"

if ! mkdir -p "$HOME_BIN_DIR"; then
  err "unable to create directory at $HOME_BIN_DIR"
fi

echo "⬇️ Download Capuchinator ${VERSION}..."
curl -L -o capuchinator.tar.gz "$URL"

echo "📦 Unpacking the archive..."
tar -xzf capuchinator.tar.gz $HOME_BIN_DIR/capuchinator

echo "🔧 Setting execution rights..."
chmod +x $HOME_BIN_DIR/capuchinator

echo "🧹 Deleting temporary files..."
rm capuchinator.tar.gz

echo "✅ Capuchinator ${VERSION} successfully installed!"
echo "🚀 Launch: capuchinator"

need_cmd() {
  if ! check_cmd "$1"; then
    err "need '$1' (command not found)"
  fi
}

check_cmd() {
  command -v "$1" > /dev/null 2>&1

  return $?
}

err() {
  local red
  local reset

  red=$(tput setaf 1 2>/dev/null || echo '')
  reset=$(tput sgr0 2>/dev/null || echo '')

  echo "${red}ERROR${reset}: $1" >&2

  exit 1
}
