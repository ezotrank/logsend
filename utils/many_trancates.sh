#!/bin/bash


i=0
while true; do
  if [ $i -le 500 ]; then
    cat $1 >> $2
  else
    echo truncate
    cat $1 > $2
    i=0
  fi
  i=$(($i+1))
done