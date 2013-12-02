#!/bin/bash

for i in out/* ; do
  cat $i >> results.csv
done
