#!/bin/bash

if [ ! -f "preferences.json" ]; then
   echo "preferences.json file not found."
   exit 1
fi

total_items=$(jq '. | length' preferences.json)
true_count=$(jq '[.[] | select(.public == true)] | length' preferences.json)
false_count=$(jq '[.[] | select(.public == false)] | length' preferences.json)

echo "Total items: $total_items"
echo "Count of true: $true_count"
echo "Count of false: $false_count"
