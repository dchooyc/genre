#!/bin/bash
# formats jsons
# requires jq

for file in ./jsons/*; do
    for jsn in $file/*.json; do
        jq . "$jsn" > temp && mv temp "$jsn"
    done
    echo "$file formatted"
done

jq . "genres.json" > temp && mv temp "genres.json"
