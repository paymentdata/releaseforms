#!/bin/bash

# given a file, pass the path to this script.
#$ cat ~/releasehistory
#epoch:fromSHA:toSHA
#1545917656:6769f55:4094e17
#1545931758:4094e17:d069bd7
#1546116375:d069bd7:885be84
#1546452309:885be84:8d9e98d

for RELEASE_SPAN in $(awk -F: '{print $2".."$3}' $1); do
    echo -e "\tgenerating csv for changelog in releasespan: ${RELEASE_SPAN}"

    commits=$(git log "${RELEASE_SPAN}" --pretty=format:"%h" --no-merges)
    commitcount=$(wc -l <<< "$commits")

    for idx in $(seq 1 ${commitcount}); do
        commit_id=$(sed "${idx}q;d" <<< "$commits")
        git log -1 --pretty=format:%at,%h,%an,%s "${commit_id}"
    done

    echo -e "\tfinished generating csv for changelog\n"
done