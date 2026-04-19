#! /usr/bin/env bash

# Benchmark the mandelbrot set project and measure it's perf across various flags and options

now=$(date "+%Y-%m-%d_%H%M%S" )

echo "Enter the benchmark name"
read benchmark_name

git checkout master

git add *

if ! git diff --cached --quiet; then
    git commit -m "Benchmarking Commit: $now"
fi

go build .

short_hash=$(git rev-parse --short HEAD)
echo "Last Commit Hash: $short_hash"

echo "List all boolean flags you want to test" 

mapfile -t boolFlags

time_data=$(/usr/bin/time -f "%e,%U,%S" ./ mandelbroset-go "${boolFlags[@]}"2>$1 > /dev/null )