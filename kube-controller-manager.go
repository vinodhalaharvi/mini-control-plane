package main

import (
	"context"
	"log"
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

	go s.MyResourceWatcher(cli)
	s.PodWatcher(cli)
}

func (s *KubeControllerManager) PodWatcher(cli *clientv3.Client) {
	watchChan := cli.Watch(context.Background(), "/pods/", clientv3.WithPrefix())
	for wresp := range watchChan {
		for _, ev := range wresp.Events {
			log.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}

func (s *KubeControllerManager) MyResourceWatcher(cli *clientv3.Client) {
	myResourcesChan := cli.Watch(context.Background(), "/myresources/", clientv3.WithPrefix())
	for response := range myResourcesChan {
		for _, ev := range response.Events {
			log.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			// Here you would add logic to reconcile the resource based on the change
		}
	}
}
