#!/bin/bash

# Script to show diff, add and commit files one by one
# Usage: ./git-diff-add-commit.sh

# Get list of modified/untracked files (excluding deleted files)
files=$(git status --porcelain | grep -v "^D" | cut -c4-)

if [ -z "$files" ]; then
    echo "No files to commit."
    exit 0
fi

echo "Found files to process:"
echo "$files"
echo ""

# Process each file
for file in $files; do
    if [ -f "$file" ]; then
        echo "=== Processing: $file ==="
        
        # Show diff (for new files, show content)
        if git ls-files --error-unmatch "$file" >/dev/null 2>&1; then
            echo "Changes in $file:"
            git diff "$file"
        else
            echo "New file $file:"
            if [ -s "$file" ]; then
                head -50 "$file"
                echo "... (showing first 50 lines)"
            else
                echo "(empty file)"
            fi
        fi
        
        echo ""
        read -p "Enter commit message for $file (or press Enter to skip): " message
        
        if [ -n "$message" ]; then
            git add "$file"
            git commit -m "$message"
            echo "✓ Committed $file"
        else
            echo "✗ Skipped $file"
        fi
        
        echo "----------------------------------------"
        echo ""
    fi
done

echo "All files processed."
