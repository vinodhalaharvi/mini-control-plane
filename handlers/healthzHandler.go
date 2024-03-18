package handlers

import (
	"github.com/gin-gonic/gin"
	"go.etcd.io/etcd/client/v3"
)

func NewHealthZHandler(client *clientv3.Client) *HealthZHandler {
	return &HealthZHandler{Client: client}
}

func (s *HealthZHandler) Healthz(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

type HealthZHandler struct {
	Client *clientv3.Client
}
