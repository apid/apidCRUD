#! /bin/bash
# identify functions that are not preceded by a comment.

badfunc()
{
	local file=$1
	awk < "$file" \
		-v "FNAME=$file" '
		BEGIN {
			hascom = 0
		}
		!hascom && /^func[ 	]/ {
			print FNAME ":" $0
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

for f in $(srcfiles *.go); do
	badfunc "$f"
done
