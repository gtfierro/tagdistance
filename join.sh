#!/bin/bash

rm -rf out
mkdir out
rm results.csv
for i in out/* ; do
  cat $i >> results.csv
done
