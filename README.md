# ipfs-cluster-erasure-example

Base on https://github.com/loomts/ipfs-cluster, provide some use cases and test the performance.

## use ali cloud machines deployment

Base on [ansible script](https://github.com/hsanjuan/ansible-ipfs-cluster) and make a little change to deploy
IPFS Cluster(EC) and test it on cloud server.

### reference

https://ipfscluster.io/documentation/reference/configuration/#:~:text=Manual%20identity%20generation

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