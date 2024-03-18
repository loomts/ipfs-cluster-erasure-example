all: ipfs ipfs-cluster
ipfs:
	ansible-playbook -i inventory.yml ipfs.yml
ipfs-cluster:
	ansible-playbook -i inventory.yml ipfs-cluster.yml
clean:
	ansible-playbook -i inventory.yml clean.yml
stop:
	ansible-playbook -i inventory.yml stop.yml
.PHONY = all ipfs ipfs-cluster clean stop