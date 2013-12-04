#!/usr/bin/env python

import sys
from collections import defaultdict

filename = sys.argv[1]
N = int(sys.argv[2])

print filename, N

numbers = set()
d = defaultdict(list)

with open(filename) as f:
    for line in f:
        tmp = line.split(',')
        d[tmp[0].strip()].append( (tmp[1].strip(), float(tmp[2].strip())) )

for key in d.iterkeys():
    l = d[key]
    l.sort(reverse=True,key=lambda x: x[1])
    strings = [','.join([x[0], str(x[1])]) for x in l[:N]]
    for s in strings:
	print key+',', s
