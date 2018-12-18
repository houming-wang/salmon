package election

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"salmon/pkg/config"
	"salmon/pkg/utils"

	"go.etcd.io/etcd/clientv3"
)

type ElectionInterface interface {
	// Start election
	Start()

	// Stop election
	Stop()

	// True if this election become leader. This value will change as leader changes
	IsLeader() bool

	// Return the current leader name
	GetCurrentLeader() string
}

type ElectionInfo struct {
	conf       *config.EtcdConfig
	client     *clientv3.Client
	leaseId    clientv3.LeaseID
	stop       chan bool // a chan for sending signal to stop the election
	name       string
	isLeader   bool
	leaderName string
}

func NewElection(name string) (*ElectionInfo, error) {
	cli, err := utils.NewEtcdV3Client()
	if err != nil {
		log.Fatal("Error while initializing connection to Etcd: ", err)
		return nil, err
	}

	etcdConf := config.GetConfig()
	ctx, cancel := context.WithTimeout(context.Background(), etcdConf.Timeout)
	defer cancel()
	// In etcd v3, lease is used to associate a k/v and set a ttl to it.
	resp, err := cli.Grant(ctx, etcdConf.LeaderTTL)
	if err != nil {
		log.Fatal("Error while granting etcd lease: ", err)
		return nil, err
	}

	return &ElectionInfo{
		conf:     etcdConf,
		client:   cli,
		leaseId:  resp.ID,
		stop:     make(chan bool),
		name:     name,
		isLeader: false,
	}, err
}

func (ei *ElectionInfo) Start() {
	// A ticker chan used to refresh key lease every ttl/4
	tickChan := time.NewTicker(time.Second * time.Duration(ei.conf.LeaderTTL/4)).C
	ei.watchLeader()
	for {
		select {
		case <-ei.stop:
			log.Println("Stop signal received.")
			return
		case <-tickChan:
			// Try to acquire the leader role and send heartbeat to renew lease
			ei.acquire()
			ei.keepAliaveOnce()
		default:
		}
	}
}

func (ei *ElectionInfo) Stop() {
	if ei.stop != nil {
		ei.stop <- true
		close(ei.stop)
	}
}

func (ei *ElectionInfo) IsLeader() bool {
	return ei.isLeader
}

func (ei *ElectionInfo) GetCurrentLeader() string {
	return ei.leaderName
}

// Watch the key changes and set the current leader info
func (ei *ElectionInfo) watchLeader() {
	watchChan := ei.client.Watch(context.Background(), ei.conf.LeaderKey)
	go func() {
		for resp := range watchChan {
			for _, ev := range resp.Events {
				v := string(ev.Kv.Value[:])
				if clientv3.EventTypePut == ev.Type && v == ei.name {
					ei.isLeader = true
					ei.leaderName = v
					log.Println(ei.name, "become leader.")
				} else if clientv3.EventTypePut == ev.Type && v != ei.name {
					ei.isLeader = false
					ei.leaderName = v
					log.Println(ei.name, "is follower.")
				}
			}
		}
	}()
}

// Use the etcd transaction(Txn) to put a key/value.
func (ei *ElectionInfo) acquire() error {
	ctx, cancel := context.WithTimeout(context.Background(), ei.conf.Timeout)
	defer cancel()
	// Use Txn instead of Put will make sure data consitent in different election client.
	txnResp, err := ei.client.Txn(ctx).
		// CreateRevision = 0 means this key is not existed
		If(clientv3.Compare(clientv3.CreateRevision(ei.conf.LeaderKey), "=", 0)).
		// Put k/v if leader key is not existed
		Then(clientv3.OpPut(ei.conf.LeaderKey, ei.name, clientv3.WithLease(ei.leaseId))).
		Else(clientv3.OpGet(ei.conf.LeaderKey)).
		Commit()
	if err != nil {
		log.Println("Error while putting election key:", err)
		return err
	}
	if !txnResp.Succeeded {
		getResp := (*clientv3.GetResponse)(txnResp.Responses[0].GetResponseRange())
		if len(getResp.Kvs) == 0 {
			errmsg := fmt.Sprintf("Leader key is not existed")
			log.Println(errmsg)
			return errors.New(errmsg)
		}
		v := string(getResp.Kvs[0].Value)
		ei.leaderName = v
		if ei.name == v {
			ei.isLeader = true
		} else {
			ei.isLeader = false
		}
	}
	return nil
}

// Keep alive once will send a heartbeat to renew the etcd lease.
func (ei *ElectionInfo) keepAliaveOnce() (*clientv3.LeaseKeepAliveResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ei.conf.Timeout)
	defer cancel()
	resp, err := ei.client.KeepAliveOnce(ctx, ei.leaseId)
	if err != nil {
		log.Println("Error while send heartbeat:", err)
		return nil, err
	}
	return resp, nil
}
