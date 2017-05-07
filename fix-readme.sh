#! /bin/bash
# replace the Template section of the README.md,
# with the current version of the template.

README=${1:-README.md}
TEMPLATE=${2:-template.txt}
PROGNAME=${0##*/}
tmpf=/tmp/$PROGNAME-$$.tmp

awk < "$README" > "$tmpf" \
	-v "TEMPLATE=$TEMPLATE" \
	'BEGIN {
		found = 0
		skip = 0
	}
	!found && /^Template:/ {
		found = 1
		skip = 1
		print $0
		print "```"
		system("cat " TEMPLATE)
		print "```"
		next
	}
	skip && /^```/ {
		skip++
		if (skip >= 3) skip = 0
		next
	}
	skip {
		next
	}
	{ print $0 }
	'

/bin/mv "$tmpf" "$README"
