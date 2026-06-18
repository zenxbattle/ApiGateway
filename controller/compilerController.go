package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"xcode/natsclient"

	"github.com/gin-gonic/gin"
)

type CompilerController struct {
	NatsClient *natsclient.Client
}

func NewCompilerController(nc *natsclient.Client) *CompilerController {
	return &CompilerController{NatsClient: nc}
}

type ExecutionRequest struct {
	Code     string `json:"code" binding:"required"`
	Language string `json:"language" binding:"required"`
}

type ExecutionResponse struct {
	Output        string `json:"output"`
	Error         string `json:"error,omitempty"`
	StatusMessage string `json:"status_message"`
	Success       bool   `json:"success"`
	ExecutionTime string `json:"execution_time,omitempty"`
}

func (s *CompilerController) CompileCodeHandler(c *gin.Context) {
	var req ExecutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ExecutionResponse{
			Error:         err.Error(),
			StatusMessage: "API request failed",
			Success:       false,
		})
		return
	}

	// Marshal the request to JSON
	reqData, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ExecutionResponse{
			Error:         err.Error(),
			StatusMessage: "Failed to process request",
			Success:       false,
		})
		return
	}

	// Send request to NATS
	msg, err := s.NatsClient.Request("compiler.execute.request", reqData, 15*time.Second)
	if err != nil {
		c.JSON(http.StatusBadGateway, ExecutionResponse{
			Error:         "Failed to execute code",
			StatusMessage: "Failed to execute code",
			Success:       false,
		})
		return
	}

	// Unmarshal the response
	var resp ExecutionResponse
	if err := json.Unmarshal(msg, &resp); err != nil {
		c.JSON(http.StatusBadRequest, ExecutionResponse{
			Error:         err.Error(),
			StatusMessage: "Failed to parse response",
			Success:       false,
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
