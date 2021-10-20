#!/bin/bash

rm -f main && go build . && strip -s ./main && ./main
