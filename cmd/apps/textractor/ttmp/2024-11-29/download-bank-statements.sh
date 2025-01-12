#!/usr/bin/env bash

# Read from file if provided, otherwise from stdin
if [ $# -eq 1 ]; then
    input=$(cat "$1")
else
    input=$(cat)
fi

mkdir -p outputs/json
mkdir -p outputs/csv

# Process each JSON object using jq
echo "$input" | jq -c '.[]' | while read -r job; do
    job_id=$(echo "$job" | jq -r '.job_id')
    document=$(echo "$job" | jq -r '.document')
    basename=$(basename "$document" .pdf)
    
    echo "Processing $basename..."
    textractor --config textractor-config.json fetch "$job_id" > "outputs/json/$basename.json"
    textractor debug csv -o "outputs/csv/$basename.csv" "outputs/json/$basename.json"
done

# run like:
# textractor --config textractor-config.json list --since today --status COMPLETED --output json | ./ttmp/2024-11-29/download-bank-statements.sh