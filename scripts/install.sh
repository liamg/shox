#!/bin/bash

set -e

echo "Determining platform..."
platform=$(uname | tr '[:upper:]' '[:lower:]')
echo "Finding latest release..."
asset=$(curl --silent https://api.github.com/repos/liamg/shox/releases/latest | grep -o "https://github.com/liamg/shox/releases/download/.*/shox-$platform-amd64" | head -n1)
echo "Downloading latest release for your platform..."
curl -s -L -H "Accept: application/octet-stream" "${asset}" --output ./shox
echo "Installing shox..."
chmod +x ./shox
installdir="${HOME}/bin/"
if [ "$EUID" -eq 0 ]; then
  installdir="/usr/local/bin/"
fi
mkdir -p $installdir
mv ./shox "${installdir}/shox"
which shox &> /dev/null || (echo "Please add ${installdir} to your PATH to complete installation!" && exit 1)
echo "Installation complete!"
