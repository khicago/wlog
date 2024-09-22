#!/bin/bash

# set output file name
OUTPUT_FILE="benchmark.txt"

# get current date and time
CURRENT_DATE=$(date "+%Y-%m-%d %H:%M:%S")

# get current git commit hash
GIT_COMMIT=$(git rev-parse HEAD)

# get device information
OS=$(uname -s)
ARCH=$(uname -m)
CPU_INFO=$(sysctl -n machdep.cpu.brand_string 2>/dev/null || lscpu | grep "Model name" | sed -r 's/Model name:\s{1,}//g' 2>/dev/null || echo "Unknown")
TOTAL_MEMORY=$(sysctl -n hw.memsize 2>/dev/null | awk '{ printf "%.2f GB", $1/1024/1024/1024 }' || free -h | awk '/^Mem:/ {print $2}' 2>/dev/null || echo "Unknown")

# run benchmark test, output log to console, and capture results
BENCHMARK_RESULTS=$(go test -bench=. -benchmem 2>&1 | grep -E "Benchmark|ns/op|allocs/op")

# output results and device information to file
{
    echo "Benchmark Results"
    echo "Date: $CURRENT_DATE"
    echo "Git Commit: $GIT_COMMIT"
    echo ""
    echo "Device Information:"
    echo "OS: $OS"
    echo "Architecture: $ARCH"
    echo "CPU: $CPU_INFO"
    echo "Total Memory: $TOTAL_MEMORY"
    echo ""
    echo "Benchmark Output:"
    echo "$BENCHMARK_RESULTS"
} > "$OUTPUT_FILE"

echo "Benchmark results have been saved to $OUTPUT_FILE"