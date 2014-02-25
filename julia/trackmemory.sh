#!/bin/bash

rm mem.log
filename=$1
julia jaccard_as_matrix.jl $filename &
pid=$!
echo $pid
while true 
  do 
    ps -p $pid -o %cpu=,%mem= >> mem.log
    sleep 2
done
