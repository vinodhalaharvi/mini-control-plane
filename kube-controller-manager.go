package main

import (
	"log"
	"mini-control-plane/watchers"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type KubeControllerManager struct {
}

func NewKubeControllerManager() *KubeControllerManager {
	return &KubeControllerManager{}
}

func (s *KubeControllerManager) Run() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to connect to etcd: %v", err)
	}
	defer cli.Close()

	log.Println("Connected to etcd successfully")

	watcher := watchers.NewDockerWatcher()
	watcher.WatchDockerPaths(cli)
}
