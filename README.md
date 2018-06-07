# bench-23
A tiny program for benchmarking concurrent file io in Go and Node.js

## Results
Tested on CentOS 7, Xeon E3 4 Cores 8 Threads.

|                                         | Go 1.10.2                              | Node.js v10.3.0                   |
|-----------------------------------------|----------------------------------------|-----------------------------------|
| 10000 files (Write, Read)               | Write: 311.389699ms Read: 313.480085ms | Write: 453.411ms Read: 871.045ms  |
| 5000 files (Write, Read and parse JSON) | Write: 154.625183ms Read: 3.489775076s | Write: 202.956ms Read: 3967.285ms |

## Run benchmarks
Go:
```sh
# build the program
GOOS=linux GOARCH=amd64 go build -o ./bin/bench23

# 10000 files
./bench23 10000

# 5000 files + parse json
./bench23 5000 --parse-json
```

Node.js (assume yarn already installed):
```sh
# 10000 files
NODE_ENV=production yarn start 10000

# 5000 files + parse json
NODE_ENV=production yarn start 5000 --parse-json
```
