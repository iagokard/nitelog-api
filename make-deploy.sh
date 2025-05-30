#!/bin/bash

ZIP_FILE="deploy.zip"

# Check for swag
command -v swag &> /dev/null || {
    echo "Error: swag not found. Install with:"
    echo "go install github.com/swaggo/swag/cmd/swag@latest"
    exit 1
}

# Handle existing zip file
if [[ -f "$ZIP_FILE" ]]; then
    read -rp "$ZIP_FILE exists. Remove? [y/N] " answer
    [[ "$answer" =~ ^[Yy]$ ]] && rm -f "$ZIP_FILE" || exit 0
fi

# Generate docs
echo "Generating Swagger docs..."
swag init -g main.go || exit 1

# Create zip
echo "Creating $ZIP_FILE..."
zip -r "$ZIP_FILE" \
    discloud.config \
    docs \
    go.mod \
    go.sum \
    internal \
    main.go \
    .env \
    nitelog-1e1a5-firebase-adminsdk-3gzzr-5803258e1f.json

echo "Created $ZIP_FILE"

