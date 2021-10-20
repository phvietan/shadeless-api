#!/bin/bash

red=`tput setaf 1`
green=`tput setaf 2`
yellow=`tput setaf 3`
blue=`tput setaf 4`
magenta=`tput setaf 5`
cyan=`tput setaf 6`
reset=`tput sgr0`

echo "Running ${green}test cases${reset} on files ${blue}*_test.go:${reset}"

mkdir -p cov
go test ./... -coverprofile cov/coverage.out | sed ''/ok/s//$(printf "${green}PASS${reset}")/''
go tool cover -html=cov/coverage.out -o cov/coverage.html

echo "======== DONE ========="
echo "Serving report on: ${magenta}http://localhost:12345/coverage.html${reset}"
python3 -m http.server 12345 --directory cov

