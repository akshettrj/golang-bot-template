#!/bin/bash

new_package_name="$(echo $1 | sed 's/^\s*//;s/\s*$//;s/\s+/ /g;s/\s/_/g')"

if [ -z "$new_package_name" ]
then
    echo "Usage:"
    echo
    echo "    $0 <new_package_name>"
    exit 1
fi

echo "New package name : '$new_package_name'"
read -p "Confirmed ? [Y/n] " -n 1 -r
if [[ ! $REPLY =~ ^[yY]$ ]]
then
    echo "Aborting..."
    exit
fi

echo

find . -type f -name '*.go' -o -name 'go.*' | xargs -I'{}' sed -si 's/golang-bot-template/'$new_package_name'/g' "{}"

echo "Done..."
