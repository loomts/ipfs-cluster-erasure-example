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

func main() {
	sth := utils.NewFileHelper()
	tree := sth.GetTreeMultiReader()
	err := ECAddFaultTolerantAndRetrive(tree)
	if err != nil {
		log.Error(err)
	} else {
		log.Infof("success to add and retrieve directory %s %d", tree.Name(), tree.Size())
	}
	// create 1KB to 1GB files then pin and retrieve them when some nodes down

	fs := sth.GetRandFileMultiReader()
	for i := 0; i < len(fs); i++ {
		err := ECAddFaultTolerantAndRetrive(fs[i])
		if err != nil {
			log.Error(err)
		}
		log.Infof("success to add and retrieve ramdon file %s %d", fs[i].Name(), fs[i].Size())
	}
}

func ECAddFaultTolerantAndRetrive(f utils.ECFile) error {
	log.Infof("start docker")
	err := StartDocker()
	if err != nil {
		return err
	}
	log.Infof("add file %s", f.Name())
	ci, err := AddECFile(f)
	if err != nil {
		return err
	}
	log.Infof("close docker")
	err = CloseDocker()
	if err != nil {
		return err
	}
	log.Infof("retrieve file %s", f.Name())
	err = RetrieveECFile(ci)
	if err != nil {
		return err
	}
	log.Infof("verify file")
	err = VerifyECFile(f.Name(), ci.String())
	if err != nil {
		return err
	}
	log.Infof("all precesses done successfully")
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

func CloseDocker() error {
	cli, err := dockercli.NewClientWithOpts(dockercli.FromEnv, dockercli.WithAPIVersionNegotiation())
	defer cli.Close()
	if err != nil {
		return err
	}
	// x := rand.Intn(2) + 1
	// nodes := []int{x, x + 1}
	x := rand.Intn(7) + 1
	nodes := []int{x, x + 1, x + 2}
	for _, node := range nodes {
		cluster := "cluster" + fmt.Sprintf("%d", node)
		if err := cli.ContainerStop(context.Background(), cluster, container.StopOptions{}); err != nil {
			return err
		}
		ipfs := "ipfs" + fmt.Sprintf("%d", node)
		if err := cli.ContainerStop(context.Background(), ipfs, container.StopOptions{}); err != nil {
			return err
		}
	}
	return nil
}

func RetrieveECFile(ci api.Cid) error {
	c, err := client.NewClient()
	if err != nil {
		return err
	}
	return c.ECGet(context.Background(), ci, utils.RetrieveDir)
}

func VerifyECFile(file string, ci string) error {
	hashSource, err := utils.HashDirectory(path.Join(utils.SourceDir, file))
	if err != nil {
		return err
	}
	hashTarget, err := utils.HashDirectory(path.Join(utils.RetrieveDir, ci))
	if err != nil {
		return err
	}
	if hashSource != hashTarget {
		return fmt.Errorf("source and target are not the same")
	}
	return nil
}
