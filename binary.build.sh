#!/bin/bash

# Check if version argument is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <version>"
    exit 1
fi

VERSION=$1
BINARY="terraform-switcher"
GOOS_LIST=("darwin" "linux" "windows")
GOARCH_LIST=("386" "amd64" "arm" "arm64")
GOARM_LIST=("6" "7")

# Function to check if a GOOS/GOARCH pair is unsupported
is_unsupported() {
    local goos=$1
    local goarch=$2

    if [[ "$goos" == "darwin" && "$goarch" == "386" ]]; then
        return 0
    elif [[ "$goos" == "darwin" && "$goarch" == "arm" ]]; then
        return 0
    elif [[ "$goos" == "windows" && "$goarch" == "arm64" ]]; then
        return 0
    else
        return 1
    fi
}

# Function to build the binary
build_binary() {
    local goos=$1
    local goarch=$2
    local goarm=$3
    local output="$BINARY_v$VERSION_${goos}_${goarch}"
    
    if [ "$goarch" == "arm" ]; then
        output="${output}v$goarm"
    fi
    
    if [ "$goos" == "windows" ]; then
        output="${output}.exe"
    fi

    echo "Building $output..."
    env GOOS=$goos GOARCH=$goarch GOARM=$goarm go build -o "./bin/$output"
    
    if [ $? -ne 0 ]; then
        echo "Failed to build $output"
    else
        tar -czf "./bin/${BINARY}_v${VERSION}_${goos}_${goarch}.tar.gz" -C ./bin "$output"
        rm "./bin/$output"
    fi
}

# Create output directory
mkdir -p bin

# Loop through each combination of GOOS and GOARCH
for goos in "${GOOS_LIST[@]}"; do
    for goarch in "${GOARCH_LIST[@]}"; do
        # Skip unsupported GOOS/GOARCH pairs
        if is_unsupported $goos $goarch; then
            continue
        fi
        
        if [ "$goarch" == "arm" ]; then
            for goarm in "${GOARM_LIST[@]}"; do
                build_binary $goos $goarch $goarm
            done
        else
            build_binary $goos $goarch ""
        fi
    done
done

echo "All binaries built successfully."
