#!/bin/bash

for d in $(go list ./... | grep -v "vendor"); do
    golint -set_exit_status $d > lintreport.out
    exit_code=$?
    if [ $exit_code -ne 0 ]; then
      exit $exit_code
    fi
done
