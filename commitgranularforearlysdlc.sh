#!/bin/bash

#RELEASE_SPAN takes the form "$fromSHA..$toSHA"
RELEASE_SPAN=$1

commits=$(git log ${RELEASE_SPAN} --pretty=format:"%h" --no-merges)
commitcount=$(wc -l <<< "$commits")

for idx in $(seq 1 ${commitcount}); do
    commit_id=$(sed "${idx}q;d" <<< "$commits")
    git log -1 --pretty=format:%at,%h,%an,%s "${commit_id}"
done
