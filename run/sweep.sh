#!/bin/bash

set -e

case=$1
test -z "$case" && exit 1

file="assets/$case.json"
test -e "$file" || exit 1

name='maxLevel'
values='2 4 6 8 10'

function tune {
  cat $file | sed -e "s/\(\"$name\"\s*:\)[^,]*/\1 $value/" > ${file}_temp
  mv ${file}_temp ${file}
}

for value in $values; do
  tune $file $name $value
  make $case
done
