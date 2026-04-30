#! /usr/bin/env bash

# Benchmark the mandelbrot set project and measure it's perf across various flags and options

set -e

if ! command -v /usr/bin/time &> /dev/null; then
    echo "Error: /usr/bin/time is not installed"
    exit 1
fi

if ! command -v bc &> /dev/null; then
    echo "Error: bc is require"
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
    time_data=$( { printf "\n\n\n\n\n" | /usr/bin/time -f "%e %U %S %M" /bin/bash -c './build/mandelbrotset "${boolFlags[@]}"  1>/dev/null 2>&1;'; } 2>&1 )

    echo $time_data
    echo "run $i finished"

    read -r real_time user_time sys_time max_mem avg_mem <<< "$time_data"
    total_real_time=$(echo "$total_real_time + $real_time" | bc)
    total_sys_time=$(echo "$total_sys_time + $sys_time" | bc)
    total_user_time=$(echo "$total_user_time + $user_time" | bc)
    cat  <<EOF >> benchmark.txt
    run $i
    total time: $real_time \t cumulative time: $total_real_time
    user time: $user_time \t cumulative time: $total_user_time
    system time: $sys_time \t cumulative time: $total_sys_time
    maximum memory: $max_mem
    average memory: $avg_mem


EOF
done

average_real_time=$(echo "scale=4; $total_real_time / 4.0" | bc)
average_user_time=$(echo "scale=4; $total_user_time / 4.0" | bc)
average_sys_time=$(echo "scale=4; $total_sys_time / 4.0" | bc)
cat >> benchmark.txt <<EOF
------- Results -----------
average real time: $average_real_time
average user time: $average_user_time
average sys time: $average_sys_time

EOF

echo "Benchmarking finished"