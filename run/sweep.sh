#!/bin/bash

set -e

case=$1
if [ -z "$case" ]; then
	echo 'Error: need a case name.'
	exit 1
fi

file="assets/${case}_ext.json"
if [ ! -f "$file" ]; then
	file="assets/${case}.json"
fi
if [ ! -f "$file" ]; then
	echo 'Error: cannot find the configuration file.'
	exit 1
fi

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
