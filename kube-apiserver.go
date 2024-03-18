package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"

	"log"
	"net/http"
)

type KubeAPIServer struct {
}

func NewKubeAPIServer() *KubeAPIServer {
	return &KubeAPIServer{}
}

func (s *KubeAPIServer) Run() {
	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})

	http.HandleFunc("/myresources", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			var resource MyResource
			if err := json.NewDecoder(r.Body).Decode(&resource); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			// Convert resource to JSON to store in etcd
			jsonData, err := json.Marshal(resource)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = cli.Put(context.Background(), "/myresources/"+resource.Name, string(jsonData))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, "Resource created")
			// Handle other methods (GET, DELETE) similarly
		}
	})

	fmt.Printf("Api Server is running at http://localhost:8080\n")
	err = r.Run(":8080")
	if err != nil {
		log.Printf("Error: %v", err)
	}
}
