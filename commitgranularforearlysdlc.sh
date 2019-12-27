#!/bin/bash

#RELEASE_SPAN takes the form "$fromSHA..$toSHA"
RELEASE_SPAN=$1

commits=$(git log ${RELEASE_SPAN} --pretty=format:"%h" --no-merges)
commitcount=$(wc -l <<< "$commits") && ((commitcount++))

for idx in $(seq 1 ${commitcount}); do
    echo $idx
    commit_id=$(sed "${idx}q;d" <<< "$commits")

    git log --name-status --diff-filter="ACDMRT" -1 -U $commit_id
done
