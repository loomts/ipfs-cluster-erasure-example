for cid in $(ipfs-cluster-ctl pin ls | awk '{print $1}'); do
  ipfs-cluster-ctl pin rm $cid
done