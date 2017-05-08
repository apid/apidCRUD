#! /bin/bash
# identify functions that are not preceded by a comment.

badfunc()
{
	awk 'BEGIN {
		hascom = 0
		}
		!hascom && /^func[ 	]/ {
			print $0
			next
		}
		/^\/\// {
			hascom = 1
			next
		}
		{
			hascom = 0
			next
		}
		'
}

srcfiles()
{
	ls "$@" | grep -v '_test\.go$'
}

# ----- start of mainline code
cat $(srcfiles "$@") | badfunc
