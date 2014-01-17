#!/bin/bash

rm results.csv > /dev/null
for i in out/* ; do
  cat $i >> results.csv
done
sed -i'' -e 's/ //g' results.csv
