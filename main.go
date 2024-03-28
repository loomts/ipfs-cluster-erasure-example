package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/apenella/go-ansible/v2/pkg/adhoc"
	"github.com/apenella/go-ansible/v2/pkg/execute"
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

	// fmt.Println("ecget_diff")
	// TestECGetFileDiffSize()
	// fmt.Println("ecget_same")
	// TestECGetFileLargeSameSize()

	//fmt.Println("ecget_recovery")
	//TestECRecovery()
	utils.Analysis()
}

func AddFaultTolerantAndRetrieve(add func(f utils.ECFile) (api.Cid, error), c *client.Client, f utils.ECFile) error {
	ci, err := add(f)
	if err != nil {
		return err
	}

	nodes, err := checkStatusAndStop(f.Name)
	if err != nil {
		return err
	}
	defer func() {
		err = restartNodes(nodes)
		if err != nil {
			log.Error(err)
		}
		fmt.Println("unpin", f.Name)
		_, err = c.Unpin(context.Background(), ci)
		if err != nil {
			log.Error(err)
		}
		// local gc
		_, err := c.RepoGC(context.Background(), true)
		if err != nil {
			log.Error(err)
		}
	}()

	fmt.Println("ecget", f.Name)
	_, err = os.Stat(utils.RetrieveDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(utils.RetrieveDir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	err = c.ECGet(context.Background(), ci, utils.RetrieveDir)
	if err != nil {
		return err
	}
	err = utils.Diff(path.Join(utils.SourceDir, f.Name), path.Join(utils.RetrieveDir, ci.String()))
	if err != nil {
		return err
	}
	return err
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
	files := sth.Get1GBRandFileMultiReader()
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
	files := sth.Get1GBRandFileMultiReader()
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
	_, err = checkStatusAndStop("")
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

func checkStatusAndStop(name string) (string, error) {
	cmd := fmt.Sprintf(`ipfs-cluster-ctl status --filter pinned --sort name | grep -A1 "%s" | grep -v "\[guangzhou-00\]" | grep ">" | awk '{print$2}' | awk -F']' '{print $2}' | sort -u | head -n 4 | xargs | tr ' ' ','`, name)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	nodes := string(out)
	fmt.Printf("stopping %s\n", out)

	adhoc := &adhoc.AnsibleAdhocCmd{
		Pattern: "all",
		AdhocOptions: &adhoc.AnsibleAdhocOptions{
			Args:       "systemctl stop ipfs; systemctl stop ipfs-cluster",
			Connection: "ssh",
			Inventory:  nodes,
			ModuleName: "shell",
			Become:     true,
		},
	}

	e := execute.NewDefaultExecute(
		execute.WithCmd(adhoc),
		execute.WithWrite(NullWriter(0)),
	)
	err = e.Execute(context.TODO())
	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)
	return nodes, nil
}

func restartNodes(nodes string) error {
	fmt.Println("starting ", nodes)

	adhoc := &adhoc.AnsibleAdhocCmd{
		Pattern: "all",
		AdhocOptions: &adhoc.AnsibleAdhocOptions{
			Args:       "systemctl start ipfs; systemctl start ipfs-cluster",
			Connection: "ssh",
			Inventory:  nodes,
			ModuleName: "shell",
			Become:     true,
		},
	}

	e := execute.NewDefaultExecute(
		execute.WithCmd(adhoc),
		execute.WithWrite(NullWriter(0)),
	)
	err := e.Execute(context.TODO())
	if err != nil {
		panic(err)
	}
	time.Sleep(5 * time.Second)
	return nil
}

type NullWriter int

func (NullWriter) Write([]byte) (int, error) {
	return 0, nil
}
