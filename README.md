# ipfs-cluster-erasure-example
Base on https://github.com/loomts/ipfs-cluster, provide some use cases and test the performance.

## Usage
```zsh
# clone ipfs-cluster erasure implementation
cd $GOPATH/src
git clone https://github.com/loomts/ipfs-cluster
git clone https://github.com/loomts/ipfs-cluster-erasure-example
cd ipfs-cluster-erasure-example
go mod tidy
go run main
```