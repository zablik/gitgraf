#!/bin/bash

echo "info,files_changed,files_created,files_deleted,lines_added,lines_deleted,approvers" >> stats3.csv

# Use git log to batch retrieve commit hashes and other details, excluding merge commits.
git log --pretty=format:'%H,%ae,%an,%ad' --no-merges | while IFS=',' read commit_hash author_email author_name author_date; do
    # Use git show to batch retrieve stats and diff information.
    read lines_added lines_deleted files_changed files_created files_deleted <<< $(git show --format="" --numstat --no-commit-id $commit_hash | awk '
    BEGIN {files_changed=0; lines_added=0; lines_deleted=0;}
    {lines_added += $1; lines_deleted += $2; files_changed++;}
    END {print lines_added, lines_deleted, files_changed}')

    # Initialize counters for files.
    files_created=0
    files_deleted=0

    # Use git diff-tree to get detailed file change information for each commit.
    while IFS= read -r line; do
        case "$line" in
            A*) ((files_created++)) ;;
            D*) ((files_deleted++)) ;;
        esac
    done < <(git diff-tree --no-commit-id --name-status -r $commit_hash)
    
    # Extract approvers from commit messages in a batch.
    approvers=$(git log -1 --pretty=format:'%b' $commit_hash | grep 'Approved-by:' | sed 's/Approved-by: //g' | jq -R -s -c 'split("\n")[:-1]')

    # Combine information.
    info="$commit_hash,$author_email,$author_name,$author_date"
    
    # Output to file.
    echo "$info,$files_changed,$files_created,$files_deleted,$lines_added,$lines_deleted,$approvers" >> stats3.csv
done
