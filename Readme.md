# GCopy

This project writed in Golang, it copies files (small or larges)
With this in mind, we're going to face some challenges:
- copy small files
- copy large files
- copy files by network
- copy files by network with networking issues


### Tests
```sh
go test ./... -v
```


**TODO: add to pipeline**
```sh
GOARCH=amd64 go build -ldflags "-X 'main.Version=0.0.1' -X 'main.BuildTime=$(date)' -X 'main.Description=Copy files around'" -o build/gcopy ./cmd/gcopy/

```