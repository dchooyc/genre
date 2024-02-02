#!/bin/bash
# formats jsons
# requires jq

for file in ./jsons/*.json; do
    jq . "$file" > temp && mv temp "$file"
done

jq . "genres.json" > temp && mv temp "genres.json"
