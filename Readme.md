# GCopy

This project is written in Golang and focuses on building a robust utility for copying files of various sizes and across different environments.

**Goals and Challenges**  
The project aims to tackle several real-world scenarios:
- Copying small files efficiently
- Handling large file transfers with stability
- Transferring files over a network
- Managing file transfers under unreliable network conditions


**How to use it**
```sh
gcopy [source] [destination]
```


### Tests
```sh
go test ./... -v
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