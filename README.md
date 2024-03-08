# ipfs-cluster-erasure-example
Base on https://github.com/loomts/ipfs-cluster, provide some use cases and test the performance.

## Docker deployment
Deploy IPFS Cluster by docker. If you want to deploy on different physical machine, see this [branch](https://github.com/loomts/ipfs-cluster-erasure-example/tree/ansible).

## Usage
```zsh
cd $GOPATH/src
git clone https://github.com/loomts/ipfs-cluster
git clone https://github.com/loomts/ipfs-cluster-erasure-example
cd ipfs-cluster-erasure-example
go mod tidy
go run main
```

### reference
https://ipfscluster.io/documentation/reference/configuration