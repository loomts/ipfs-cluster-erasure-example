all: ipfs ipfs-cluster
ipfs:
	ansible-playbook -i inventory.ini ipfs.yml
ipfs-cluster:
	ansible-playbook -i inventory.ini ipfs-cluster.yml
clean-and-restart:
	ansible-playbook -i inventory.yml ipfs.yml --tags "clean and restart"
stop:
	ansible-playbook -i inventory.ini stop.yml
start:
	ansible-playbook -i inventory.ini start.yml
update-ipfs-cluster-service:
	ansible-playbook -i inventory.ini update-ipfs-cluster-service.yml

.PHONY = all ipfs ipfs-cluster clean-and-restart stop start update-ipfs-cluster-service