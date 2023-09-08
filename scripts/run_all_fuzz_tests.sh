#!/bin/bash

# Find all _fuzz_test.go files and extract their directories
for dir in $(find . -name '*_fuzz_test.go' -exec dirname {} \; | sort -u); do
    echo "Fuzzing package in directory: $dir"
    go test -fuzz=. -fuzztime=1m -v $dir
done