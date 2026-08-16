# GCopy

A CLI tool written in Go that copies files of any size — from small to over 10 GB — across local filesystems and networks. It streams data during transfer so memory usage stays low regardless of file size, and it can dynamically increase or decrease the buffer/package size mid-transfer to optimize performance based on actual conditions (e.g., slow disk, network latency).

## Getting Started

### Installation

**From source:**
```sh
go install github.com/paolinos/gcopy@latest
```

Or build from source:
```sh
make all
```

---

**Key capabilities:**

- Copy files locally between directories in parallel
- Transfer files over a TCP-based network protocol with framing
- Handle unreliable networks by retrying failed transfers and detecting dead connections
- Preserve file metadata (permissions, timestamps) during local copies
- Skip or overwrite existing destination files based on user choice



## Display Format
Show up to N files copying simultaneously, each with a progress bar and percentage:
```
source: ~/tmp/folder
destination: ~/ext
Total: 3 / 10 files copied
file1.txt       ████████████████████████████████████████████████████████████████████████████░░░    94.5 MiB  15.9 MiB/s 00:06 95%
file2.log       ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░    5120 MiB  0.00 MiB/s 00:00 0%
{FILE_NAME}     ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░    {FILE_SIZE} MiB {Transfer rate} MiB/s {Time} {Percentage}%
```
- Progress bar uses ASCII-safe characters (█ filled, ░ empty) for terminal compatibility
- Update progress every N bytes or on newline flush to prevent interleaved output
- Show total files copied vs. total files when all transfers complete

### CLI Flags
| Flag	        | Default           |	Description
| --parallel	| 5	                | Max concurrent file transfers
| --chunks      | 0 (auto)          | Number of chunks for parallel copying
| --check       | false             | Validate each file's source and destination sizes after copy
| --skip        | false             | Skip existing files at the destination without overwriting
| --force       | false             | Overwrite all existing files at the destination
| --show, -n    | 5                 | Number of files to display in results summary

### Exit Codes
0: All copies completed successfully
1: One or more copies failed (partial success acceptable)
2: Invalid arguments or configuration errors


### Usage

Copy files between directories in parallel:
```sh
gcopy [source] [destination]
```


### Tests
```sh
go test -v ./... 

go test -cover -race ./...
```


**TODO: add to pipeline**
```sh
GOARCH=amd64 go build -ldflags "-X 'main.Version=0.0.1' -X 'main.BuildTime=$(date)' -X 'main.Description=Copy files around'" -o build/gcopy ./cmd/gcopy/

```

### Build with make
Review the [Makefile](./Makefile) to see the options
```
make all
```