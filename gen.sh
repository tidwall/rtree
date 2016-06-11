#!/bin/bash

set -e

cd $(dirname "${BASH_SOURCE[0]}")

rm -rf vendor/

gen(){
	if [ "$2" == "true" ]; then
		dex="d"
	else
		dex=""
	fi
	mkdir -p vendor/d$1$dex
	cat rtree_base.go | \
		sed "s/TNUMDIMS/$1/g" | \
		sed "s/TDEBUG/$2/g" | \
		sed 's/\/\/\ +build ignore/\/\/generated; DO NOT EDIT!/g' | \
		gofmt > vendor/d$1$dex/rtree.go
}

for i in {1..4}; do 
	#gen $i true
	gen $i false
done

