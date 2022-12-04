#!/usr/bin/env bash

count=0

max1=0
max2=0
max3=0


while read -r line; do

    if [ -z "$line" ]; then
        echo "total count: $count"
        # check count with max
        if [ "$count" -gt "$max1" ]; then
            max3=$max2
            max2=$max1
            max1=$count
        elif [ "$count" -gt "$max2" ]; then
            max3=$max2
            max2=$count
        elif [ "$count" -gt "$max3" ]; then
            max3=$count
        fi

        # reset count
        count=0
        continue
    else
        ((count+=$line))
    fi
done < input.txt

echo $max1
echo $max2
echo $max3
max=$(( $max1 + $max2 + $max3 ))
echo "Max: $max"
