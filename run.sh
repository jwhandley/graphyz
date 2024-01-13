#!/bin/zsh

# Check if an argument is provided
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 filename"
    exit 1
fi

# Use the first argument as the filename
filename=$1

# Run the command with the filename
switcherooctl launch ./graphyz "$filename"
