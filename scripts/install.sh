#!/bin/bash

# Проверяем, передан ли аргумент с версией
if [ -z "$1" ]; then
  echo "💡 Usage: $0 <version>"
  echo "Example: $0 v0.1.0"
  exit 1
fi

VERSION=$1
TARBALL="capuchinator_${VERSION#v}_linux_amd64.tar.gz"
URL="https://github.com/capuchinapp/capuchinator/releases/download/${VERSION}/capuchinator_${VERSION#v}_linux_amd64.tar.gz"

echo "⬇️ Download Capuchinator version ${VERSION}..."
curl -L -o capuchinator.tar.gz "$URL"

echo "📦 Unpacking the archive..."
tar -xzf capuchinator.tar.gz capuchinator

echo "🔧 Setting execution rights..."
chmod +x ./capuchinator

echo "🧹 Deleting temporary files..."
rm capuchinator.tar.gz

echo "✅ Capuchinator ${VERSION} successfully installed!"
echo "🚀 Launch: ./capuchinator"
