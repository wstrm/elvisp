#!/usr/bin/env bash

set -e

usage=$(cat <<- EOM
Usage:\n
./test_full.bash github.com/user/pkg
EOM
)

if [ -z "$1" ]; then echo -e $usage && exit 1; fi

pkgname() {
	name=(${1//\// })
	echo ${name[2]}
}

pkg="$1"
main="$(pkgname $pkg)"
dir=$(dirname $0)

echo $main

local_dependencies() {
	suffix=$1
	go list -f '{{ join .Deps  "'${suffix}'\n"}}' $pkg | grep $main | grep -v vendor
}

go_vet() {
	go vet $(local_dependencies /...)
}

go_lint() {
	go run golang.org/x/lint/golint -min_confidence 0.0 $(local_dependencies /...)
}

go_test() {
	touch coverage.tmp
	echo 'mode: atomic' > coverage.txt && local_dependencies | xargs -n1 -I{} sh -c 'echo "> {}"; go test -mod=vendor -tags test -race -covermode=atomic -coverprofile=coverage.tmp {} && tail -n +2 coverage.tmp >> coverage.txt'
	rm coverage.tmp 2>/dev/null
}

log() {
	echo -e "\e[1m\e[34m>> \e[93m$1\e[0m"
}

log "Golint"
go_lint

log "Go vet"
go_vet

log "Go test"
go_test

log "Exit: $?"
