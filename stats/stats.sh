#!/bin/bash

# Grab the source code
go get -d github.com/$REPONAME/...

# Grab Go package name
pkgName=github.com/$REPONAME

# Grab just first path listed in GOPATH
goPath="${GOPATH%%:*}"

# Construct Go package path
pkgPath="$goPath/src/$pkgName"

if [ -e "$pkgPath/Godeps/_workspace" ];
then
  # Add local godeps dir to GOPATH
  GOPATH=$pkgPath/Godeps/_workspace:$GOPATH
fi

# get the number of dependencies in the repo
go list $pkgName/... > dep.log || true
deps=`wc -l dep.log | cut -d' ' -f1`;
rm dep.log

# get number of lines of go code
golines=`( find $pkgPath -name '*.go' -print0 | xargs -0 cat ) | wc -l`

# output the stats
echo $REPONAME, $deps, $golines
