#!/bin/bash

echo "launching comparator"

BASE="$PWD/.."
FOLDER="results/api-benchmark-official-mix"
BASE_PATH="${BASE}/${FOLDER}"
OUTPUT_FOLDER="$PWD/$1"

mkdir "$OUTPUT_FOLDER"

declare -a SETS=("1")
QUERIES="1000"

PROJECT="_ach_10sec_delay_"

echo $(which python3)

for set in ${SETS[@]}
do
    echo "analyzing set of $set"
    python3 analyzer.py $OUTPUT_FOLDER \
        ${BASE_PATH}/prysm_arch_10_secs/prysm${PROJECT}${QUERIES}_${set}-* \
        ${BASE_PATH}/lighthouse_arch_10_secs/lighthouse${PROJECT}${QUERIES}_${set}-* \
        ${BASE_PATH}/teku_arch_10_secs/teku${PROJECT}${QUERIES}_${set}-* \
        ${BASE_PATH}/nimbus_arch_10_secs/nimbus${PROJECT}${QUERIES}_${set}-* 

    mkdir "f_comparison_${set}"
    mv comparison_* "f_comparison_${set}"

done

