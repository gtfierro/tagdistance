#!/bin/bash

rm mem.log
filename=$1
julia jaccard_as_matrix.jl $filename &
pid=$!
echo $pid
while kill -0 $pid >/dev/null 2>&1
  do 
    ps -p $pid -o %cpu=,%mem= >> mem.log
    sleep 2
done
