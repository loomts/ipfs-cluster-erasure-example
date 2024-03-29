package client

import (
	"context"
	"sync"
	"time"

	"github.com/ipfs-cluster/ipfs-cluster/api"
	"github.com/ipfs-cluster/ipfs-cluster/api/rest/client"
	files "github.com/ipfs/boxo/files"
	logging "github.com/ipfs/go-log/v2"
	ma "github.com/multiformats/go-multiaddr"
)

type Client struct {
	client.Client
}

var log = logging.Logger("EClient")
var defaultTimeout = 3 * time.Minute
var ECAddParams, DefaultAddParams api.AddParams

func init() {
	ECAddParams = api.DefaultAddParams()
	ECAddParams.ReplicationFactorMin = 1
	ECAddParams.ReplicationFactorMax = 1
	ECAddParams.Erasure = true // automatically enable shard and raw-leaves
	ECAddParams.DataShards = 4
	ECAddParams.ParityShards = 2
	ECAddParams.Recursive = true

	DefaultAddParams = api.DefaultAddParams()
	DefaultAddParams.Shard = true
	DefaultAddParams.RawLeaves = true
	DefaultAddParams.ReplicationFactorMin = 3
	DefaultAddParams.ReplicationFactorMax = 3
	// when adding folders, Recursive must be true
	DefaultAddParams.Recursive = true
}

func NewClient() (*Client, error) {
	addr := "/ip4/127.0.0.1/tcp/9094"
	maddr, err := ma.NewMultiaddr(addr)
	if err != nil {
		return &Client{}, err
	}
	cfg := client.Config{APIAddr: maddr}
	c, err := client.NewLBClient(&client.Failover{}, []*client.Config{&cfg}, 1)
	if err != nil {
		return &Client{}, err
	}
	return &Client{c}, nil
}

func (c *Client) AddMultiFile(ctx context.Context, multiFileR *files.MultiFileReader, params api.AddParams) (api.AddedOutput, error) {
	ctx, cancle := context.WithTimeout(ctx, defaultTimeout)
	defer cancle()
	out := make(chan api.AddedOutput, 100)
	var added api.AddedOutput
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range out {
			if v.Name == params.Name {
				log.Infof("cluster pinned file, name: %s, cid: %s", v.Name, v.Cid)
				added = v
			}
		}
	}()
	err := c.Client.AddMultiFile(ctx, multiFileR, params, out)
	if err != nil {
		return api.AddedOutput{}, err
	}
	wg.Wait()
	return added, nil
}

func (c *Client) Add(ctx context.Context, path string, params api.AddParams) (api.AddedOutput, error) {
	ctx, cancle := context.WithTimeout(ctx, defaultTimeout)
	defer cancle()
	out := make(chan api.AddedOutput, 100)
	var added api.AddedOutput
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range out {
			if v.Name == params.Name {
				// log.Infof("cluster pinned file, name: %s, cid: %s", v.Name, v.Cid)
				added = v
			}
		}
	}()
	err := c.Client.Add(ctx, []string{path}, params, out)
	if err != nil {
		return api.AddedOutput{}, err
	}
	wg.Wait()
	return added, nil
}
