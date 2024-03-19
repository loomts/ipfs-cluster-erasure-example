package main

import (
	"context"
	"fmt"
	"os/exec"
	"path"
	"time"

	"github.com/ipfs-cluster/ipfs-cluster/api"
	logging "github.com/ipfs/go-log/v2"
	"github.com/loomts/ipfs-cluster-erasure-example/client"
	"github.com/loomts/ipfs-cluster-erasure-example/utils"
)

var log = logging.Logger("cluster")

// before test, ensure cluster is set up and running
func main() {
	// fmt.Println("get_diff")
	// TestGetFileDiffSize()
	// fmt.Println("get_same")
	// TestGetFileLargeSameSize()

	fmt.Println("ecget_diff")
	TestECGetFileDiffSize()
	// fmt.Println("ecget_same")
	// TestECGetFileLargeSameSize()

	//fmt.Println("ecget_recovery")
	//TestECRecovery()
	// utils.Draw()
}

func AddFaultTolerantAndRetrieve(add func(f utils.ECFile) (api.Cid, error), c *client.Client, f utils.ECFile) error {
	ci, err := add(f)
	if err != nil {
		return err
	}
	nodes, err := checkStatusAndStop()
	if err != nil {
		return err
	}
	defer func() {
		err = restartNodes(nodes)
		if err != nil {
			log.Error(err)
		}
		_, err = c.RepoGC(context.Background(), false)
		if err != nil {
			log.Error(err)
		}
	}()
	err = c.ECGet(context.Background(), ci, utils.RetrieveDir)
	if err != nil {
		return err
	}
	err = utils.Diff(path.Join(utils.SourceDir, f.Name), path.Join(utils.RetrieveDir, ci.String()))
	if err != nil {
		return err
	}
	return nil
}

func TestGetFileDiffSize() {
	sth := utils.NewFileHelper()
	files := sth.GetRandFileMultiReader()
	defer sth.Clean()
	c, err := client.NewClient()
	if err != nil {
		log.Error(err)
	}
	for _, f := range files {
		err = AddFaultTolerantAndRetrieve(c.AddFile, c, f)
		if err != nil {
			log.Error(err)
		}
	}
}

func TestGetFileLargeSameSize() {
	sth := utils.NewFileHelper()
	files := sth.Get512MBRandFileMultiReader()
	defer sth.Clean()
	c, err := client.NewClient()
	if err != nil {
		log.Error(err)
	}
	for _, f := range files {
		err = AddFaultTolerantAndRetrieve(c.AddFile, c, f)
		if err != nil {
			log.Error(err)
		}
	}
}

func TestECGetFileDiffSize() {
	sth := utils.NewFileHelper()
	files := sth.GetRandFileMultiReader()
	defer sth.Clean()
	c, err := client.NewClient()
	if err != nil {
		log.Error(err)
	}
	for _, f := range files {
		err = AddFaultTolerantAndRetrieve(c.AddECFile, c, f)
		if err != nil {
			log.Error(err)
		}
	}
}

func TestECGetFileLargeSameSize() {
	sth := utils.NewFileHelper()
	files := sth.Get512MBRandFileMultiReader()
	defer sth.Clean()
	c, err := client.NewClient()
	if err != nil {
		log.Error(err)
	}
	for _, f := range files {
		err = AddFaultTolerantAndRetrieve(c.AddECFile, c, f)
		if err != nil {
			log.Error(err)
		}
	}
}

func TestECRecovery() {
	sth := utils.NewFileHelper()
	files := sth.GetRandFileMultiReader()
	defer sth.Clean()
	c, err := client.NewClient()
	if err != nil {
		log.Error(err)
	}
	cis := make([]api.Cid, 0)
	for _, f := range files {
		ci, err := c.AddECFile(f)
		cis = append(cis, ci)
		if err != nil {
			log.Error(err)
		}
	}
	err = stopNodes([]string{})
	if err != nil {
		log.Error(err)
	}
	out := make(chan api.Pin, 1024)
	st := time.Now()
	err = c.ECRecovery(context.Background(), out)
	if err != nil {
		log.Error(err)
	}
	fmt.Printf("ecrecovery use %v\n", time.Since(st))
	for r := range out {
		fmt.Printf("%v\n", r)
	}
}

func checkStatusAndStop() ([]string, error) {
	cmd := exec.Command("bash", "-c", "dctl status --filter pinned --sort name")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	nodes := []string{string(output)}
	fmt.Printf("status output: %s\n", output)
	err = stopNodes(nodes)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func stopNodes(nodes []string) error {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("for node in %s; do ansible $node -m ansible.builtin.systemd -a \"name=ipfs state=stopped\" -b; ansible $node -m ansible.builtin.systemd -a \"name=ipfs-cluster state=stopped\" -b; done", nodes))
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func restartNodes(nodes []string) error {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("for node in %s; do ansible $node -m ansible.builtin.systemd -a \"name=ipfs state=started\" -b; ansible $node -m ansible.builtin.systemd -a \"name=ipfs-cluster state=started\" -b; done", nodes))
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
