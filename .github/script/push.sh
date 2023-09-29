#!/bin/sh

if [ -z "$(git status --porcelain)" ]; then
    echo "nothing to update."
else
    git add .
    git commit -m "update img.txt ${GITHUB_SHA}" -a
    git push origin main
fi
