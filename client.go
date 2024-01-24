package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ipfs-cluster/ipfs-cluster/api"
	"github.com/ipfs-cluster/ipfs-cluster/api/rest/client"
	files "github.com/ipfs/boxo/files"
	shell "github.com/ipfs/go-ipfs-api"
	logging "github.com/ipfs/go-log/v2"
	ma "github.com/multiformats/go-multiaddr"
)

type Client struct {
	client.Client
}

var log = logging.Logger("cluster")
var defaultTimeout = 30 * time.Second

func NewClient() (Client, error) {
	addr := "/ip4/127.0.0.1/tcp/9094"
	maddr, err := ma.NewMultiaddr(addr)
	if err != nil {
		return Client{}, err
	}
	cfg := client.Config{APIAddr: maddr}
	c, err := client.NewLBClient(&client.Failover{}, []*client.Config{&cfg}, 1)
	if err != nil {
		return Client{}, err
	}
	return Client{c}, nil
}

func (c *Client) AddECMultiFile(ctx context.Context, multiFileR *files.MultiFileReader, fname string) (api.AddedOutput, error) {
	addParams := api.DefaultAddParams()
	addParams.ReplicationFactorMin = 1
	addParams.ReplicationFactorMax = 1
	addParams.Name = fname
	addParams.ShardSize = 1024 * 1024

	ctx, cancle := context.WithTimeout(ctx, defaultTimeout)
	defer cancle()
	out := make(chan api.AddedOutput)
	go func() {
		c.Client.AddMultiFile(ctx, multiFileR, addParams, out)
	}()
	var added api.AddedOutput
	select {
	case added = <-out:
		log.Infof("%s added to Cluster, cid:%s\n", added.Cid)
	case <-ctx.Done():
		return api.AddedOutput{}, errors.New("timeout AddECFile")
	}

	return added, nil
}

func (c *Client) GetFileDag(ci string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	shell := c.Client.IPFS(ctx)
	printMerkleDAG(shell, ci, 0)
}

func printMerkleDAG(shell *shell.Shell, ci string, indent int) {
	var out map[string]interface{}
	err := shell.DagGet(ci, &out)
	if err != nil {
		log.Infof("failed to get dag", err)
		return
	}

	fmt.Printf("%s%s\n", strings.Repeat(" ", indent), ci)
	for k, v := range out {
		fmt.Printf("%s%s: %v\n", strings.Repeat(" ", indent+2), k, v)
	}

	links, ok := out["Links"].([]interface{})
	if ok {
		for _, link := range links {
			linkMap, ok := link.(map[string]interface{})
			if ok {
				childCid, ok := linkMap["Hash"].(string)
				if ok {
					printMerkleDAG(shell, childCid, indent+2)
				}
			}
		}
	}
}
