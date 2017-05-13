#! /usr/bin/python
#	python rtest.py N
# read the same record N times, for timing purposes.

import os
import sys
import requests
import string

HOST = "localhost"
PORT = 9000
TABLE = "bundles"
BASE_URL = "http://%s:%d/apid" % (HOST, PORT)

def apid_read(table, id):
	url = "%s/db/_table/%s/%d" % (BASE_URL, table, id)
	r = requests.get(url)
	# print "r=", r.status_code, "text=", r.text

def main(n):
	for i in range(n):
		apid_read(TABLE, 1)

if __name__ == "__main__":
	N = string.atoi(sys.argv[1])
	main(N)
