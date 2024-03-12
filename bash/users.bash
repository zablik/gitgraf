#!/bin/bash

# Check if correct arguments are provided
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 /path/to/repo branch_name"
    exit 1
fi

REPO_PATH="$1"
BRANCH_NAME="$2"

# Check if the directory exists
if [ ! -d "$REPO_PATH" ]; then
    echo "The specified directory does not exist: $REPO_PATH"
    exit 1
fi

# Navigate to the repository directory
cd "$REPO_PATH" || exit

# Check if we're in a Git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "This is not a Git repository: $REPO_PATH"
    exit 1
fi

# Switch to the desired branch
git checkout "$BRANCH_NAME" &> /dev/null

# Determine the script's parent directory for reports
SCRIPT_DIR=$(dirname "$(realpath "$0")")
PARENT_DIR=$(dirname "$SCRIPT_DIR")
REPO_NAME=$(basename "$REPO_PATH")

# Directory for the CSV file, within "reports" directory in the script's parent directory
CSV_DIR="${PARENT_DIR}/reports/${REPO_NAME}"
mkdir -p "$CSV_DIR"

# Prepare the CSV file
CSV_PATH="${CSV_DIR}/users.csv"
echo "email,name,number of commits" > "$CSV_PATH"

# Generate CSV entries
declare -A email_names

# Process each commit to populate or update the email-names array
git log --pretty="%ae,%an" | while read -r email name; do
    if [ -z "${email_names[$email]}" ]; then
        email_names[$email]="$name"
    fi
done

# Count commits per email, using the most recent name stored in the associative array
git log --pretty="%ae" | sort | uniq -c | sort -b -k1,1nr | while read -r count email; do
    name=${email_names[$email]}
    # Append data to CSV
    echo "\"$email\",\"$name\",$count" >> "$CSV_PATH"
done

echo "CSV file generated: $CSV_PATH"
