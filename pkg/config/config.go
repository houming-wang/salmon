package config

import (
	"flag"
	"time"

	"github.com/spf13/pflag"
)

var (
	endpoints = pflag.String("endpoints", "", "The addresses of etcd cluster")
	cert      = pflag.String("cert", "", "identify secure client using this TLS certificate file")
	key       = pflag.String("key", "", "identify secure client using this TLS key file")
	cacert    = pflag.String("cacert", "", "verify certificates of TLS-enabled secure servers using this CA bundle")
	timeout   = pflag.Duration("timeout", 3*time.Second, "timeout(format: 3s) for etcd dial and operation")
	leaderkey = pflag.String("leaderkey", "/easystack/rd/leader", "key in etcd for election")
	leaderttl = pflag.Int64("leaderttl", 6, "time-to-live in seconds for etcd key/value")
)

var conf *EtcdConfig

type EtcdConfig struct {
	Endpoints string        `json:"endpoints"`
	Cert      string        `json:"cert"`
	Key       string        `json:"key"`
	CACert    string        `json:"cacert"`
	Timeout   time.Duration `json:"timeout"`
	LeaderKey string        `json:"leaderkey"`
	LeaderTTL int64         `json:"leaderttl"`
}

func init() {
	conf = loadConfig()
}

func loadConfig() *EtcdConfig {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	return &EtcdConfig{
		Endpoints: *endpoints,
		Cert:      *cert,
		Key:       *key,
		CACert:    *cacert,
		Timeout:   *timeout,
		LeaderKey: *leaderkey,
		LeaderTTL: *leaderttl,
	}
}

func GetConfig() *EtcdConfig {
	return conf
}
