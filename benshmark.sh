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

total_real_time=0
total_user_time=0
total_sys_time=0

cat  <<EOF >> benchmark.txt
==============================================================================================================================
name: $benchmark_name
commit hash: $short_hash
flags: "${boolFlags[@]}"
EOF

for i in {1..4}; do
    time_data=$( { printf "\n\n\n\n\n" | /usr/bin/time -f "%e %U %S" /bin/bash -c './build/mandelbrotset "${boolFlags[@]}"  1>/dev/null 2>&1;'; } 2>&1 )

    echo $time_data
    echo "Benchmarking finished"

    read -r real_time user_time sys_time <<< "$time_data"
    ((total_real_time += real_time))
    ((total_user_time += user_time))
    ((total_sys_time += sys_time))
    cat  <<EOF >> benchmark.txt
    run $i
    total time: $real_time
    user time: $user_time
    system time: $sys_time


EOF
done

cat >> benchmark.txt <<EOF
------- Results -----------
average real time: ((total_real_time / 3))
average user time: ((total_user_time / 3))
average sys time: ((total_sys_time / 3))

EOF
