#!/usr/bin/env bash

max=0
count=0

while read -r line; do

    if [ -z "$line" ]; then
        echo "total count: $count"
        # check count with max
        if [ "$count" -gt "$max" ]; then
            echo "new max: $count"
            max=$count
        fi
        # reset count
        count=0
        continue
    else
        ((count+=$line))
    fi
done < input.txt

echo "Max: $max"
