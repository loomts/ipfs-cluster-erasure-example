# cd $GOPATH/src && git clone github.com/loomts/ipfs-cluster
cd $GOPATH/src/ipfs-cluster/cmd/ipfs-cluster-ctl && make
cd $GOPATH/src/ipfs-cluster
docker build -t ipfs-cluster-erasure -f $GOPATH/src/ipfs-cluster/Dockerfile-erasure .
sleep 2
cd $GOPATH/src/ipfs-cluster-erasure-example
docker-compose -f ipfs-cluster-erasure.yml up -d