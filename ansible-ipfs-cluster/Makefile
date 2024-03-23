all: ipfs ipfs-cluster
ipfs:
	ansible-playbook -i inventory.yml ipfs.yml
ipfs-cluster:
	ansible-playbook -i inventory.ini ipfs-cluster.yml
clean-and-restart:
	ansible-playbook -i inventory.ini clean.yml
	ansible-playbook -i inventory.ini ipfs.yml --tags "clean and restart"
	ansible-playbook -i inventory.ini ipfs-cluster.yml --tags "clean and restart"
stop:
	ansible-playbook -i inventory.ini stop.yml
start:
	ansible-playbook -i inventory.ini start.yml
clean:
	ansible-playbook -i inventory.ini clean.yml
update-ipfs-cluster-service:
	ansible-playbook -i inventory.ini update-ipfs-cluster-service.yml

.PHONY = all ipfs ipfs-cluster clean-and-restart clean stop start update-ipfs-cluster-service