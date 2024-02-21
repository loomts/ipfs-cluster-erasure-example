package main

import (
	"context"
	"fmt"
	"math/rand"
	"os/exec"
	"path"
	"time"

	"github.com/ipfs-cluster/ipfs-cluster/api"
	logging "github.com/ipfs/go-log/v2"
	"github.com/loomts/ipfs-cluster-erasure-example/client"
	"github.com/loomts/ipfs-cluster-erasure-example/utils"

	"github.com/docker/docker/api/types/container"
	dockercli "github.com/docker/docker/client"
)

var log = logging.Logger("cluster")

// case default timeout of ipfs-cluster shard retrieve is 3 minutes, this main function need 11(file number)*3 minutes at most.
func main() {
	log.Infof("start docker compose")
	err := StartDocker()
	if err != nil {
		log.Error(err)
	}
	c, err := client.NewClient()
	if err != nil {
		log.Error(err)
	}
	sth := utils.NewFileHelper()
	tree := sth.GetTreeMultiReader()
	err = ECAddFaultTolerantAndRetrive(c, tree)
	if err != nil {
		log.Error(err)
	} else {
		log.Infof("success to add and retrieve directory %s %d bit", tree.Name(), tree.Size())
	}

	// create 1KB to 1GB files then pin and retrieve them when some nodes down
	fs := sth.GetRandFileMultiReader()
	for i := 0; i < len(fs); i++ {
		err := ECAddFaultTolerantAndRetrive(c, fs[i])
		if err != nil {
			log.Error(err)
			continue
		}
		log.Infof("success to add and retrieve ramdon file %s %d bit", fs[i].Name(), fs[i].Size())
	}
}

func ECAddFaultTolerantAndRetrive(c client.Client, f utils.ECFile) error {
	log.Infof("add file %s", f.Name())
	ci, err := AddECFile(f)
	if err != nil {
		return err
	}
	nodes, err := CloseContainers()
	if err != nil {
		return err
	}
	defer func() {
		log.Infof("restart containers")
		err = ReStartContainers(nodes)
		if err != nil {
			log.Error(err)
		}
		log.Infof("ipfs gc")
		_, err = c.RepoGC(context.Background(), false)
		if err != nil {
			log.Error(err)
		}
	}()
	log.Infof("retrieve file %s", f.Name())
	err = c.ECGet(context.Background(), ci, utils.RetrieveDir)
	if err != nil {
		return fmt.Errorf("%s ERROR:%s", ci, err)
	}
	log.Infof("verify file")
	err = utils.Diff(path.Join(utils.SourceDir, f.Name()), path.Join(utils.RetrieveDir, ci.String()))
	if err != nil {
		return err
	}
	log.Infof("%s successfully ECAddFaultTolerantAndRetrive", f.Name())
	return nil
}

func StartDocker() error {
	cmd := exec.Command("/bin/zsh", "start.sh")
	err := cmd.Run()
	time.Sleep(10 * time.Second) // wait for cluster peers set up
	return err
}

func AddECFile(f utils.ECFile) (api.Cid, error) {
	defer f.Closer.Close()
	c, err := client.NewClient()
	if err != nil {
		return api.CidUndef, err
	}
	ctx := context.Background()
	params := client.ECAddParams
	params.Name = f.Name()
	pin, err := c.AddMultiFile(ctx, f.Mfr, params)
	if err != nil {
		return api.CidUndef, err
	}
	return pin.Cid, nil
}

func CloseContainers() ([]int, error) {
	cli, err := dockercli.NewClientWithOpts(dockercli.FromEnv, dockercli.WithAPIVersionNegotiation())
	defer cli.Close()
	if err != nil {
		return nil, err
	}
	x := rand.Intn(6) + 1
	nodes := []int{x, x + 1, x + 2, x + 3}
	for _, node := range nodes {
		cluster := "cluster" + fmt.Sprintf("%d", node)
		if err := cli.ContainerStop(context.Background(), cluster, container.StopOptions{}); err != nil {
			return nil, err
		}
		ipfs := "ipfs" + fmt.Sprintf("%d", node)
		if err := cli.ContainerStop(context.Background(), ipfs, container.StopOptions{}); err != nil {
			return nil, err
		}
	}
	log.Infof("close %v containers", nodes)
	return nodes, nil
}

func ReStartContainers(nodes []int) error {
	cli, err := dockercli.NewClientWithOpts(dockercli.FromEnv, dockercli.WithAPIVersionNegotiation())
	defer cli.Close()
	if err != nil {
		return err
	}
	for _, node := range nodes {
		cluster := "cluster" + fmt.Sprintf("%d", node)
		if err := cli.ContainerStart(context.Background(), cluster, container.StartOptions{}); err != nil {
			return err
		}
		ipfs := "ipfs" + fmt.Sprintf("%d", node)
		if err := cli.ContainerStart(context.Background(), ipfs, container.StartOptions{}); err != nil {
			return err
		}
	}
	time.Sleep(10 * time.Second)
	log.Infof("start %v containers", nodes)
	return nil
}
