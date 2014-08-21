#!/bin/bash

ROOTPATH=$(readlink -f "$(dirname $0)/..")
cachefile='/tmp/.gnotify-build.cache'

check_status() {
  local RED=$(tput setaf 1)
  local GREEN=$(tput setaf 2)
  local RESET=$(tput sgr0)

  if [[ "$1" != "0" ]]; then
    echo " [ ${RED}FAILED$RESET ]"
    cat $cachefile
    finish $1
  fi

  echo " [ ${GREEN}OK$RESET ]"
  cat $cachefile
}

finish() {
  # Mop up the temporary file
  rm -f $cachefile
  exit $1
}

printf 'Installing dependencies from requirements.txt...'
$ROOTPATH/scripts/dependencies.sh &> $cachefile
check_status $?

printf 'Compiling source...'
GOPATH="$ROOTPATH" go install xantoria.com &> $cachefile
check_status $?

finish 0
