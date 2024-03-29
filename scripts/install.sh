#!/bin/bash

install_prefix='/usr/local'

if [[ "$1" != "" ]]; then
  install_prefix="$(realpath $1)" || exit 1
fi

echo "This script will install binaries, config, etc. to the system."
echo "It'll use the install prefix '$install_prefix'."
echo "Make sure you're running this as a user with write permissions there."
echo -n 'OK? '
read -n 1 resp
echo ''

if [[ "$resp" != "y" ]]; then
  echo 'Aborting.'
  exit 1
fi

# Make sure the directory structure is present
mkdir -p "$install_prefix/bin"
mkdir -p "$install_prefix/etc"
mkdir -p "$install_prefix/share"

# Add init script
cp -f 'scripts/init.d/gnotify' '/etc/init.d/gnotify'

# Copy binary, conf, and shared files
cp -f 'bin/gnotify' "$install_prefix/bin/gnotify"
cp -n 'etc/gnotify.conf' "$install_prefix/etc/gnotify.conf"
cp -R 'share' "$install_prefix"

checksum=$(sha1sum "$install_prefix/bin/gnotify" | cut -d ' ' -f 1)
echo "Binary $install_prefix/bin/gnotify, SHA1{$checksum} installed"
