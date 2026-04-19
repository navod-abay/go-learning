#! /usr/bin/env bash

# Benchmark the mandelbrot set project and measure it's perf across various flags and options

set -e

if ! command -v /usr/bin/time &> /dev/null; then
    echo "Error: /usr/bin/time is not installed"
    exit 1
fi

now=$(date "+%Y-%m-%d_%H%M%S" )

echo "Enter the benchmark name"
read benchmark_name

git checkout master

git add .

if ! git diff --cached --quiet; then
    git commit -m "Benchmarking Commit: $now"
fi

go build -o build/mandelbrotset .

short_hash=$(git rev-parse --short HEAD)
echo "Last Commit Hash: $short_hash"

echo "List all boolean flags you want to test" 

mapfile -t boolFlags || true

echo "Finished reading command line arguments"

time_data=$( { /usr/bin/time -f "%e %U %S" ./build/mandelbrotset "${boolFlags[@]}" > /dev/null 2>&1 <<<EOF
\n
\n
\n
\n
\n
EOF
} )

echo $time_data
echo "Benchmarking finished"

read -r real_time user_time sys_time <<< "$time_data"

cat  <<EOF >> benchmark.txt
==============================================================================================================================
name: $benchmark_name
commit hash: $short_hash
flags: "${boolFlags[@]}"
total time: $real_time
user time: $user_time
system time: $sys_time
EOF