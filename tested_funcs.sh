#! /bin/bash

cat *_test.go \
| sed -n -e 's/^func Test_\([^()]*\).*/\1/p' \
| sed -e p -e 's/_[^_]*$//' \
| sort -u
