package main

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/docker/docker/api/types/container"
	dockercli "github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/ipfs-cluster/ipfs-cluster/api"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
	"github.com/loomts/ipfs-cluster-erasure-example/client"
	"github.com/loomts/ipfs-cluster-erasure-example/utils"
)

var log = logging.Logger("cluster")

func main() {
	// TestAllAndFaultToler()

	// fmt.Println("ecadd_diff")
	// TestAddECFile_DiffSize()
	// fmt.Println("ecadd_same")
	// TestAddECFile_LargeSameSize()

	// fmt.Println("add_diff")
	// TestAddFile_DiffSize()

	// fmt.Println("add_same")
	// TestAddFile_LargeSameSize()

	// fmt.Println("get_diff")
	// TestGetFile_DiffSize()
	// fmt.Println("get_same")
	// TestGetFile_LargeSameSize()
	// fmt.Println("ecget_diff")
	// TestECGetFile_DiffSize()
	// fmt.Println("ecget_same")
	// TestECGetFile_LargeSameSize()

	// fmt.Println("ecget_recovery")
	// TestECRecovery()
	// utils.Draw()

	httpserver()
}

func httpserver() {
	r := gin.Default()

	r.GET("/ecget", func(c *gin.Context) {
		ci := c.Query("cid")
		client, err := client.NewClient()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		decodedCid, err := cid.Decode(ci)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		_, err = os.Stat(utils.RetrieveDir)
		if os.IsNotExist(err) {
			err = os.Mkdir(utils.RetrieveDir, os.ModePerm)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		}
		err = client.ECGet(context.Background(), api.NewCid(decodedCid), utils.RetrieveDir)
		out := make(chan api.GlobalPinInfo, 1024)
		client.StatusCids(context.Background(), []api.Cid{api.NewCid(decodedCid)}, false, out)
		name := ""
		for r := range out {
			name = r.Name
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		encodedFilename := url.QueryEscape(name)
		c.Writer.Header().Set("Content-Disposition", "attachment; filename="+encodedFilename)
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
		c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		file, err := os.Open(utils.RetrieveDir + "/" + ci)
		if err != nil {
			log.Error("Error opening file: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		defer file.Close()

		if _, err := io.Copy(c.Writer, file); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	})
	r.Run(":8888")
}

func start() {
	cmd := exec.Command("bash", "-c", "./start.sh")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	time.Sleep(150)
}

func stopAndRemoveDockerContainers() {
	cmd := exec.Command("bash", "-c", "docker stop $(docker ps -a -q) && docker rm $(docker ps -a -q) && sudo rm -rf compose")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	time.Sleep(30)
}

// case default timeout of ipfs-cluster shard retrieve is 3 minutes, this main function need 11(file number)*3 minutes at most.
func TestAllAndFaultToler() {
	begin := time.Now()
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
	err = AddFaultTolerantAndRetrieve(AddECFile, c, tree)
	if err != nil {
		log.Error(err)
	} else {
		log.Infof("success to add and retrieve directory %s %d bit", tree.Name, tree.Size)
	}

	// create 1KB to 1GB files then pin and retrieve them when some nodes down
	fs := sth.GetRandFileMultiReader()
	for i := 0; i < len(fs); i++ {
		err := AddFaultTolerantAndRetrieve(AddECFile, c, fs[i])
		if err != nil {
			log.Error(err)
			continue
		}
		log.Infof("success to add and retrieve ramdon file %s %d bit", fs[i].Name, fs[i].Size)
	}
	log.Infof("total time: %v", time.Since(begin))
	sth.Clean()
}

func AddFaultTolerantAndRetrieve(add func(f utils.ECFile) (api.Cid, error), c *client.Client, f utils.ECFile) error {
	ci, err := add(f)
	if err != nil {
		return err
	}
	nodes, err := CloseContainers()
	if err != nil {
		return err
	}
	defer func() {
		err = ReStartContainers(nodes)
		if err != nil {
			log.Error(err)
		}
		_, err = c.RepoGC(context.Background(), false)
		if err != nil {
			log.Error(err)
		}
	}()
	start := time.Now()

	fmt.Println("ecget", f.Name)
	_, err = os.Stat(utils.RetrieveDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(utils.RetrieveDir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	// alternative method is utils.IPFSGet(ci.Cid)
	err = c.ECGet(context.Background(), ci, utils.RetrieveDir)
	if err != nil {
		return fmt.Errorf("%s ERROR:%s", ci, err)
	}
	fmt.Printf("retrieve %s %d use %v\n", ci, f.Size, time.Since(start))
	err = utils.Diff(path.Join(utils.SourceDir, f.Name), path.Join(utils.RetrieveDir, ci.String()))
	if err != nil {
		return err
	}
	// log.Infof("%s successfully AddFaultTolerantAndRetrive", f.Name())
	return nil
}

func StartDocker() error {
	cmd := exec.Command("/bin/zsh", "start.sh")
	err := cmd.Run()
	time.Sleep(10 * time.Second) // wait for cluster peers set up
	return err
}

func AddFile(f utils.ECFile) (api.Cid, error) {
	defer f.Closer.Close()
	c, err := client.NewClient()
	if err != nil {
		return api.CidUndef, err
	}
	ctx := context.Background()
	params := client.DefaultAddParams
	params.Name = f.Name
	// st := time.Now()
	pin, err := c.Add(ctx, path.Join(f.Base, f.Name), params)
	if err != nil {
		return api.CidUndef, err
	}
	// fmt.Printf("add %s(%v), size:%v, use %v\n", pin.Name, pin.Cid, f.Size(), time.Since(st))
	return pin.Cid, nil
}

func AddECFile(f utils.ECFile) (api.Cid, error) {
	defer f.Closer.Close()
	c, err := client.NewClient()
	if err != nil {
		return api.CidUndef, err
	}
	ctx := context.Background()
	params := client.ECAddParams
	params.Name = f.Name
	st := time.Now()
	pin, err := c.Add(ctx, path.Join(f.Base, f.Name), params)
	if err != nil {
		return api.CidUndef, err
	}
	fmt.Printf("add %s(%v), size:%v, use %v\n", pin.Name, pin.Cid, f.Size, time.Since(st))
	return pin.Cid, nil
}

func CloseContainers() ([]int, error) {
	cli, err := dockercli.NewClientWithOpts(dockercli.FromEnv, dockercli.WithAPIVersionNegotiation())
	defer cli.Close()
	if err != nil {
		return nil, err
	}
	// cluster[0~5]
	x := rand.Intn(4) + 1
	nodes := []int{x, x + 1}
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
	time.Sleep(10 * time.Second) // metric timeout is 30s, after 30s, peer id could remove from the cluster
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

func TestAddECFile_DiffSize() {
	// start()
	// defer func() {
	// 	stopAndRemoveDockerContainers()
	// }()
	sth := utils.NewFileHelper()
	files := sth.GetRandFileMultiReader()
	// defer sth.Clean()
	for _, f := range files {
		_, err := AddECFile(f)
		if err != nil {
			log.Error(err)
		}
	}
}

func TestAddECFile_LargeSameSize() {
	start()
	defer func() {
		stopAndRemoveDockerContainers()
	}()
	sth := utils.NewFileHelper()
	files := sth.Get512MBRandFileMultiReader()
	defer sth.Clean()
	for _, f := range files {
		_, err := AddECFile(f)
		if err != nil {
			log.Error(err)
		}
	}
}

func TestAddFile_DiffSize() {
	// start()
	// defer func() {
	// 	stopAndRemoveDockerContainers()
	// }()
	sth := utils.NewFileHelper()
	files := sth.GetRandFileMultiReader()
	// defer sth.Clean()
	for _, f := range files {
		_, err := AddFile(f)
		if err != nil {
			log.Error(err)
		}
	}
}

func TestAddFile_LargeSameSize() {
	start()
	defer func() {
		stopAndRemoveDockerContainers()
	}()
	sth := utils.NewFileHelper()
	files := sth.Get512MBRandFileMultiReader()
	defer sth.Clean()
	for _, f := range files {
		_, err := AddFile(f)
		if err != nil {
			log.Error(err)
		}
	}
}

func TestGetFile_DiffSize() {
	start()
	defer func() {
		stopAndRemoveDockerContainers()
	}()
	sth := utils.NewFileHelper()
	files := sth.GetRandFileMultiReader()
	defer sth.Clean()
	c, err := client.NewClient()
	if err != nil {
		log.Error(err)
	}
	for _, f := range files {
		err = AddFaultTolerantAndRetrieve(AddFile, c, f)
		if err != nil {
			log.Error(err)
		}
	}
}

func TestGetFile_LargeSameSize() {
	start()
	defer func() {
		stopAndRemoveDockerContainers()
	}()
	sth := utils.NewFileHelper()
	files := sth.Get512MBRandFileMultiReader()
	defer sth.Clean()
	c, err := client.NewClient()
	if err != nil {
		log.Error(err)
	}
	for _, f := range files {
		err = AddFaultTolerantAndRetrieve(AddFile, c, f)
		if err != nil {
			log.Error(err)
		}
	}
}

func TestECGetFile_DiffSize() {
	start()
	defer func() {
		stopAndRemoveDockerContainers()
	}()
	sth := utils.NewFileHelper()
	files := sth.GetRandFileMultiReader()
	defer sth.Clean()
	c, err := client.NewClient()
	if err != nil {
		log.Error(err)
	}
	for _, f := range files {
		err = AddFaultTolerantAndRetrieve(AddECFile, c, f)
		if err != nil {
			log.Error(err)
		}
	}
}

func TestECGetFile_LargeSameSize() {
	start()
	defer func() {
		stopAndRemoveDockerContainers()
	}()
	sth := utils.NewFileHelper()
	files := sth.Get512MBRandFileMultiReader()
	defer sth.Clean()
	c, err := client.NewClient()
	if err != nil {
		log.Error(err)
	}
	for _, f := range files {
		err = AddFaultTolerantAndRetrieve(AddECFile, c, f)
		if err != nil {
			log.Error(err)
		}
	}
}

func TestECRecovery() {
	start()
	defer func() {
		stopAndRemoveDockerContainers()
	}()
	sth := utils.NewFileHelper()
	files := sth.GetRandFileMultiReader()
	defer sth.Clean()
	c, err := client.NewClient()
	if err != nil {
		log.Error(err)
	}
	cis := make([]api.Cid, 0)
	for _, f := range files {
		ci, err := AddECFile(f)
		cis = append(cis, ci)
		if err != nil {
			log.Error(err)
		}
	}
	_, err = CloseContainers()
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
