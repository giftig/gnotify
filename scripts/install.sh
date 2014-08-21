#!/bin/bash

install_prefix='/usr/local/'

if [[ "$1" != "" ]]; then
  install_prefix="$1"
fi

echo "This script will install binaries, config, etc. to the system."
echo "It'll use the install prefix '$install_prefix'."
echo "Make sure you're running this as a user with write permissions there."
read "OK? " resp

if [[ "$(grep -E '^y')" != "" ]]; then
  echo 'Aborting.'
  exit 1
fi

# Make sure the directory structure is present
mkdir -p "$install_prefix/bin"
mkdir -p "$install_prefix/etc"
mkdir -p "$install_prefix/share"

# Copy binary, conf, and shared files
cp -f 'bin/xantoria.com' "$install_prefix/bin/gnotify"
cp -f 'etc/gnotify.conf' "$install_prefix/etc/gnotify.conf"
cp -R 'share/*' "$install_prefix/share"
