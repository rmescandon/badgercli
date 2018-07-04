#!/bin/sh

set -eu

if ! which dep >/dev/null;then
    export PATH="$PATH:${GOPATH%%:*}/bin"

    if ! which dep >/dev/null;then
	    echo Installing dep
	    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
        dep init
    fi
fi

echo Ensure dependencies
dep ensure -v
