package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"mini-control-plane/handlers"
)

type KubeAPIServer struct {
	Client *clientv3.Client
	Routes *gin.Engine
}

func NewKubeAPIServer() *KubeAPIServer {
	return &KubeAPIServer{}
}

func (s *KubeAPIServer) setupRoutes() {
	s.Routes = gin.Default()
	healthZHandler := handlers.NewHealthZHandler(s.Client)
	dockerHandler, err := handlers.NewDockerHandler()
	if err != nil {
		log.Fatalf("Failed to create Docker handler: %v", err)
	}

	s.Routes.POST("/containers/create", dockerHandler.CreateContainer)
	s.Routes.POST("/containers/:id/start", dockerHandler.StartContainer)
	s.Routes.POST("/containers/:id/stop", dockerHandler.StopContainer)

	s.Routes.GET("/healthz", healthZHandler.Healthz)
}

func (s *KubeAPIServer) Run() {
	fmt.Printf("Api Server is running at http://localhost:8080\n")
	err := s.Routes.Run(":8080")
	if err != nil {
		log.Printf("Error: %v", err)
	}
}
