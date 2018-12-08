package utils

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"crypto/tls"

	"salmon/pkg/config"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/pkg/transport"
)

func NewEtcdV3Client() (*clientv3.Client, error) {
	conf := config.GetConfig()
	if len(conf.Endpoints) == 0 {
		errmsg := fmt.Sprintf("--endpoints must be specified.")
		log.Fatalf(errmsg)
		return nil, errors.New(errmsg)
	}
	var tlsConf *tls.Config
	if len(conf.Cert) == 0 || len(conf.Key) == 0 || len(conf.CACert) == 0 {
		tlsConf = nil
	} else {
		tlsInfo := transport.TLSInfo{
			CertFile: conf.Cert,
			KeyFile:  conf.Key,
			CAFile:   conf.CACert,
		}
		var err error
		tlsConf, err = tlsInfo.ClientConfig()
		if err != nil {
			log.Println("Error while creating TLS config:", err)
			return nil, err
		}
	}

	endpoints := strings.Split(conf.Endpoints, ",")
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		TLS:         tlsConf,
		DialTimeout: conf.Timeout,
	}

	etcdClient, err := clientv3.New(cfg)
	if err != nil {
		log.Println("Error while connecting to etcd:", err)
		return nil, err
	}
	log.Println("Successfully connected to etcd:", conf.Endpoints)
	return etcdClient, nil
}
