# Chapter 5: Improving the Performance of CLI Tools

## Overview

- It is important to ensure that CLI tools perform well and efficiently, especially when creating tools that will deal with a large amount of information, like data analysis tools.
- In this chapter, we make a CLI application that executes statistical operations on a CSV file.
  - The program will receive two optional input parameters:
    1. `-col`: The column on which to execute the operation (defaults to 1)
    2. `-op`: The operation to execute on the selected column, either sum or average.
- Here, we make use of go's benchmarking and profiling tools to optimize our tool and reduce the overall time it takes to run.

## Benchmarked Improvements Before and After Adding Concurrency
```
benchcmp benchmarking_results/benchresults00m.txt benchmarking_results/benchresults02m.txt
benchcmp is deprecated in favor of benchstat: https://pkg.go.dev/golang.org/x/perf/cmd/benchstat
benchmark          old ns/op     new ns/op     delta
BenchmarkRun-8     613107604     89566283      -85.39%

benchmark          old allocs     new allocs     delta
BenchmarkRun-8     5043294        2531627        -49.80%

benchmark          old bytes     new bytes     delta
BenchmarkRun-8     535700634     215145695     -59.84%
```

## Directory Structure
```
.
├── csv.go      // functions to process csv data
├── errors.go   // error definitions
└── main.go     // main() and run() functions
```

