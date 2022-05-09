#!/bin/bash

echo "launching Api-Benchmark"

CLI_NAME="api-benchmark"

# Benchmark Values
BM_NAME="prysm_grpc_arch_10sec_delay"
HOST_ENDP="localhost:4000"
QUERY="/eth/v1/beacon/states/{beacon_state_number}/validator_balances?id={validator_id}"
REPLACES='["{beacon_state_number}", "{validator_id}"]'
RANGE_VALUES='["0:3526300", "0:21063"]'
QUERY_FILE="base_1M_queries.txt"
QUERY_BACKUP="import"
NUM_QUERIES=1000
SET_QUERY_DELAY=10 # Seconds

declare -a CONCURRENT_RATIOS  
CONCURRENT_RATIOS=(1)
echo "Concurrent ratios for the test: ${CONCURRENT_RATIOS[*]}"

CONFIG_FILE_FOLDER="config-files"

# $1 - bm_name
compose_conf_file()
{
    conf_file="${CONFIG_FILE_FOLDER}/${1}_conf.json"
    echo "$conf_file"
}

# $1 - base_name
# $2 - num-queries
# $3 - concurrent value
compose_bm_name()
{
    comp_name="${1}_${2}_${3}"
    echo "$comp_name"
}


# $1 - bm_name
# S2 - config_file_name
# S3 - concurrent_ratio value
compose_conffile()
{
    conf_file=$2
    echo "composing config file for $conf_file"
    # compose the configfile
    echo '{' > $conf_file
    echo '    "benchmark-name":' "\"${1}\"," >> $conf_file
    echo '    "host-endpoint":' "\"$HOST_ENDP\"," >> $conf_file
    echo '    "query":' "\"$QUERY\"," >> $conf_file
    echo '    "replaces":' "$REPLACES," >> $conf_file
    echo '    "range-values":' "$RANGE_VALUES," >> $conf_file
    echo '    "query-file":' "\"$QUERY_FILE\"," >> $conf_file
    echo '    "query-backup":' "\"$QUERY_BACKUP\"," >> $conf_file
    echo '    "num-queries":' "$NUM_QUERIES," >> $conf_file
    echo '    "concurrent-req":' "$3," >> $conf_file
    echo '    "set-query-delay":' "$SET_QUERY_DELAY" >> $conf_file
    echo '}' >> $conf_file

}


# go build the tool always before launching it 
go build -o $CLI_NAME

# main loop for the tests
for idx in ${CONCURRENT_RATIOS[@]}
do
    bm_name=$(compose_bm_name $BM_NAME $NUM_QUERIES $idx )
    echo "launching benchmark $bm_name"
    conf_file=$(compose_conf_file $bm_name)
    echo "config file at $conf_file"

    compose_conffile $bm_name $conf_file $idx

    "./$CLI_NAME" run --config-file="$conf_file"

done

