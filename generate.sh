#!/bin/bash

# Find all subdirectories containing a Go file
directories=$(find . -type f -name '*.go' -exec dirname {} \; | sort -u)

# Loop through each directory and run 'go generate'
for dir in $directories; do
  echo "Running 'go generate' in $dir"
  (cd "$dir" && go generate)
done
