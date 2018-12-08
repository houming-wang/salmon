# Salmon - the base encapsulated Go library/framework code for quickly build applications runs on top of distributed system.

**Note**: The repo is in early stage and under frequent development.

This repo will contains the following commonly used libraries:
* *Election*: framework for leader election. Early Stage
* *K8sClient*: TBD
* *OpenStackClient*: TBD 
* *RestfulAPI*: TBD


## Requirements
- go version 1.11+: this repo use go moudule(feature begins in go 1.11) to manage dependencies.
- docker version > 1.13

## Getting started

### Set up local etcd cluster or use existing etcd cluster
[Set up a local etcd cluster]https://github.com/etcd-io/etcd/blob/master/Documentation/dev-guide/local_cluster.md

### Build the project

```bash
# Linux env
$ make build
# or MacOS env
$ make build-mac
```

### Test it
Start two instance of example

```bash
$ bin/salmon --endpoints http://127.0.0.1:2379,http://127.0.0.1:22379,http://127.0.0.1:32379
```
