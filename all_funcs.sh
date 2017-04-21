#! /bin/bash

nontestfiles()
{
	ls *.go | grep -v '_test.go$'
}

grepfuncs()
{
	sed -n -e 's/^func *([^()]*) */func /' -e 's/^func *\([^()]*\).*/\1/p'
}

# ----- start of mainline
cat $(nontestfiles) | grepfuncs | sort
