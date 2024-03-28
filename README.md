# ipfs-cluster-erasure-example
Base on https://github.com/loomts/ipfs-cluster, provide some use cases and test the performance.

## Cloud server deployment
Base on [ansible script](https://github.com/hsanjuan/ansible-ipfs-cluster) and make a little change to deploy
IPFS Cluster(EC) and test it on cloud servers.

### Step
1. Make sure you setup golang envirment then clone this project and [ipfs-cluster(EC)](https://github.com/loomts/ipfs-cluster)
```zsh
cd $GOPATH/src
git clone https://github.com/loomts/ipfs-cluster
git clone https://github.com/loomts/ipfs-cluster-erasure-example
cd ipfs-cluster-erasure-example
git checkout ansible
```
2. Make sure you have some cloud server and already set ssh public key.

3. Follow [ansible-ipfs-cluster](ansible-ipfs-cluster/README.md)

## Test result
Use 37 ali cloud servers, see [result](result).

| Method                  | Average Throughput(MB/s) |
| ----------------------- | ------------------------ |
| add --erasure           | 7.57                     |
| ecget(some shards loss) | 7.43                     |
| ecrecovery              | 10.33                    |

### reference
https://ipfscluster.io/documentation/reference/configuration
