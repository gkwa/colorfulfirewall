#!/bin/bash

if [ ! -f "preferences.json" ]; then
    echo "preferences.json file not found."
    exit 1
fi

total_items=$(jq '. | length' preferences.json)
public_count=$(jq '[.[] | select(.public == true)] | length' preferences.json)
private_count=$(jq '[.[] | select(.public == false)] | length' preferences.json)

private=$(jq '[.[] | select(.public == false)] | .[] | .imagePath' preferences.json | xargs du -ch | tail -n 1 | awk '{print $1}')
public=$(jq '[.[] | select(.public == true)] | .[] | .imagePath' preferences.json | xargs du -ch | tail -n 1 | awk '{print $1}')

echo "$total_items images total: $private_count private ($private), $public_count public ($public)"
