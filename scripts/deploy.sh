#!/bin/bash

ROOTPATH=$(readlink -f "$(dirname $0)/..")
PACKAGE_NAME='gnotify'
tempdir=$(mktemp -d --tmpdir=/tmp/ ${PACKAGE_NAME}.XXXXXXXX)
cachefile="/tmp/.$PACKAGE_NAME-deploy.cache"
dump_location="giftiger_wunsch@giftig-1:/tmp/$PACKAGE_NAME-release.tar.gz"

touch $cachefile

# Build the application first
$ROOTPATH/scripts/build.sh || exit 1

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
  # Mop up the temporary files
  rm -f $cachefile &> /dev/null
  rm -Rf $tempdir &> /dev/null
  rm "$PACKAGE_NAME.tar.gz" &> /dev/null
  exit $1
}

cd $ROOTPATH

printf "Creating package $PACKAGE_NAME.tar.gz..."
mkdir "$tempdir/$PACKAGE_NAME"
cp -R --preserve=all bin "$tempdir/$PACKAGE_NAME"
cp -R --preserve=all etc "$tempdir/$PACKAGE_NAME"
cp -R --preserve=all share "$tempdir/$PACKAGE_NAME"
cp -R --preserve=all scripts "$tempdir/$PACKAGE_NAME"
cp -R --preserve=all *.md "$tempdir/$PACKAGE_NAME"
tar -pczf "${PACKAGE_NAME}.tar.gz" -C $tempdir "$PACKAGE_NAME"  &> /dev/null
check_status $?

scp "${PACKAGE_NAME}.tar.gz" "$dump_location" &> $cachefile
status=$?
printf "Copying package to $dump_location..."
check_status $status

echo ''
echo "Now log in to giftig-1 and run the install script!"
finish 0
