#!/bin/bash

#RELEASE_SPAN takes the form "$fromSHA..$toSHA"
RELEASE_SPAN=$1

commitcount=$(git log ${RELEASE_SPAN}  --oneline | grep -iv merge | wc -l) && ((commitcount++))

for idx in $(seq 1 ${commitcount}); do
    echo $idx
    git log --name-status --diff-filter="ACDMRT" -1 -U $(git log ${RELEASE_SPAN}  --oneline | grep -iv merge | sed "${idx}q;d" -- | cut -d' ' -f1)
done
