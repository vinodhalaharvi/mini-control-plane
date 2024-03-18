package watchers

import (
	"context"
	"go.etcd.io/etcd/client/v3"
	"log"
)

type DockerWatcher struct {
}

func NewDockerWatcher() *DockerWatcher {
	return &DockerWatcher{}
}

func (d *DockerWatcher) WatchDockerPaths(cli *clientv3.Client) {
	log.Println("Starting to watch Docker paths...")
	dockerChan := cli.Watch(context.Background(), "/docker/containers/", clientv3.WithPrefix())
	for response := range dockerChan {
		for _, ev := range response.Events {
			log.Printf("Docker Event - Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}
