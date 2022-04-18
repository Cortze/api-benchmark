#!/bin/bash

FILE_NAME="lodestar.txt"

RAW=$(cat $FILENAME)

IFS=',' read -ra VALIDATORS <<< "$RAW"

for validator in "${VAIDATORS[@]}";
do 
    echo "$i"
done

done <<< "$RAW"