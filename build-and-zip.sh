#!/bin/bash

# Set build environment variables
export GOARCH="arm64"
export GOOS="linux"

# Define the root directory for functions and output directory
functionsRoot="functions"
outputDir="bin"

# Ensure the output directory exists
mkdir -p $outputDir

# Loop through each folder in the functions directory
for dir in "$functionsRoot"/*/; do
  folderName=$(basename "$dir")
  mainFilePath="$dir/main.go"

  # Check if main.go exists in the folder
  if [[ -f "$mainFilePath" ]]; then
    outputBinary="$outputDir/bootstrap"
    outputZip="$outputDir/$folderName.zip"

    # Build the Go binary
    echo "Building $mainFilePath..."
    go build -ldflags="-s -w" -o "$outputBinary" "$mainFilePath"

    # Compress the binary into a ZIP file
    echo "Creating ZIP: $outputZip"
    zip -j "$outputZip" "$outputBinary"

    # Clean up the binary after zipping
    rm -f "$outputBinary"
  else
    echo "Skipped folder $folderName: main.go not found"
  fi
done
