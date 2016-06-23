#!/bin/bash
set -e #exit on an error

function WARNING(){
    echo -e "\e[101m\e[97m[WARNING]\e[49m\e[39m $@"
}

function INFO(){
    echo -e "\e[104m\e[97m[INFO]\e[49m\e[39m $@"
}

function exists() {
    type $1 > /dev/null 2>&1
}

WARNING "Usually you don't need to run this script."
WARNING "You can install ntimes by just running \"go get github.com/AkihiroSuda/ntimes\"."

exists go || (INFO "Please install go: https://golang.org/dl/"; false)
exists gometalinter || (INFO "Please install gometalinter: https://github.com/alecthomas/gometalinter"; false)
exists go-md2man || (INFO "Please install go-md2man (v1.0.5): https://github.com/cpuguy83/go-md2man"; false)

INFO "Building"
set -x
go build
go-md2man -in ntimes.1.md -out ntimes.1
set +x

INFO "Testing"
set -x
go test -v -race -cover
gometalinter --deadline 100s
set +x

INFO "Done!"
