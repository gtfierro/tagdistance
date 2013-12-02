#!/usr/bin/env python

import sys
from collections import defaultdict

filename = sys.argv[1]
N = int(sys.argv[2])

print filename, N

numbers = set()

with open(filename) as f:
    for line in f:
        numbers.add(line.split(' ')[0])

for number in numbers:
    compares = []
    with open(filename) as f:
        for line in f:
            tmp = line.split(' ')
            if tmp[0] == number:
                compares.append((tmp[1],tmp[2].strip()))
        compares.sort(key=lambda x: x[1])
        strings = [' '.join(x) for x in compares[:N]]
        print number+',', ','.join(strings)
        
