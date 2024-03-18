package handlers

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"go.etcd.io/etcd/client/v3"
	"net/http"
	"time"
)

type DockerHandler struct {
	Client     *client.Client
	EtcdClient *clientv3.Client
}

func NewDockerHandler() (*DockerHandler, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	etcdCli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})

	return &DockerHandler{
		Client:     cli,
		EtcdClient: etcdCli,
	}, nil
}

type CreateContainerRequest struct {
	Image string   `json:"image"`
	Cmd   []string `json:"cmd"`
}

func (d *DockerHandler) CreateContainer(c *gin.Context) {
	var req CreateContainerRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use the context package for context management
	ctx := context.Background()

	// Create the container with the provided image and command from the request
	resp, err := d.Client.ContainerCreate(ctx, &container.Config{
		Image: req.Image, // Use the image specified in the request
		Cmd:   req.Cmd,   // Use the command specified in the request
	}, nil, nil, nil, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//write it to etcd
	//Convert resource to JSON to store in etcd
	jsonData, err := json.Marshal(resp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_, err = d.EtcdClient.Put(context.Background(), "/docker/containers/"+resp.ID, string(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Container created", "container_id": resp.ID})
}

func (d *DockerHandler) StartContainer(c *gin.Context) {
	containerID := c.Param("id")
	ctx := context.Background()
	if err := d.Client.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// write to etcd
	// Convert resource to JSON to store in etcd
	startJson := map[string]bool{"started": true}
	jsonData, err := json.Marshal(startJson)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_, err = d.EtcdClient.Put(context.Background(), "/docker/containers/"+containerID, string(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, gin.H{"message": "Container started", "container_id": containerID})
}

func (d *DockerHandler) StopContainer(c *gin.Context) {
	containerID := c.Param("id")
	ctx := context.Background()
	if err := d.Client.ContainerStop(ctx, containerID, container.StopOptions{}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	stopJson := map[string]bool{"stopped": true}
	jsonData, err := json.Marshal(stopJson)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = d.EtcdClient.Put(context.Background(), "/docker/containers/"+containerID, string(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Container stopped", "container_id": containerID})
}
