#!/bin/bash

OUT_FILE_NAME="prysm_rewards.csv"
INIT_SLOT=144751
FINAL_SLOT=181984
VALIDATOR_INDEX_FILE="validator_indexes/prysm.txt"

./api-benchmark rewards --outfile=$OUT_FILE_NAME --init-slot=$INIT_SLOT --final-slot=$FINAL_SLOT --validator-indexes=$VALIDATOR_INDEX_FILE