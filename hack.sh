#!/bin/bash
set -e #exit on an error

function ERROR(){
    echo -e "\e[101m\e[97m[ERROR]\e[49m\e[39m $@"
}

function WARNING(){
    echo -e "\e[101m\e[97m[WARNING]\e[49m\e[39m $@"
}

function INFO(){
    echo -e "\e[104m\e[97m[INFO]\e[49m\e[39m $@"
}

function exists() {
    type $1 > /dev/null 2>&1
}

function init() {
    WARNING "Usually you don't need to run this script."
    WARNING "You can install ntimes by just running \"go get github.com/AkihiroSuda/ntimes\"."
    exists go || (ERROR "Please install go: https://golang.org/dl/"; false)
}

# target
function build(){
    exists go-md2man || (ERROR "Please install go-md2man (v1.0.5): https://github.com/cpuguy83/go-md2man"; false)
    INFO "Building"
    set -x
    go build
    go-md2man -in ntimes.1.md -out ntimes.1
    set +x
}

# target
function test(){
    exists gometalinter || (ERROR "Please install gometalinter: https://github.com/alecthomas/gometalinter"; false)
    INFO "Testing"
    set -x
    go test -i
    go test -v -race -cover
    gometalinter --deadline 100s
    set +x
}

# target
function release(){
    INFO "Building release binaries"
    set -x
    GOOS=darwin GOARCH=amd64 go build -o ntimes-darwin-x86_64
    GOOS=linux GOARCH=amd64 go build -o ntimes-linux-x86_64
    GOOS=windows GOARCH=amd64 go build -o ntimes-windows-x86_64.exe
    set +x
}


init
TARGETS="build test release"
[ $# -eq 0 ] && (ERROR "No target specified. available targets: $TARGETS"; false)
for f in $@; do
    if [ $f = "all" ]; then
	for g in $TARGETS; do $g; done
    elif [[ $TARGETS =~ $f ]]; then
	$f
    else
	ERROR "Unknown target $f. available targets: $TARGETS"
	false
    fi
done
