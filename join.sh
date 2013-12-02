#!/bin/bash

rm results.csv
for i in out/* ; do
  cat $i >> results.csv
done
