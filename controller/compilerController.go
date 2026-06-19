package controller

import (
	"encoding/json"
	"net/http"
	"strings"
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

type rateLimitResponse struct {
	Success       bool   `json:"success"`
	Error         string `json:"error"`
	StatusMessage string `json:"status_message"`
	ExecutionTime string `json:"execution_time"`
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

	reqData, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ExecutionResponse{
			Error:         err.Error(),
			StatusMessage: "Failed to process request",
			Success:       false,
		})
		return
	}

	clientIP := c.GetHeader("X-Forwarded-For")
	if clientIP != "" {
		clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
	}

	headers := map[string]string{}
	if clientIP != "" {
		headers["X-Client-IP"] = clientIP
	}

	msg, err := s.NatsClient.RequestMsg("compiler.execute.request", reqData, headers, 15*time.Second)
	if err != nil {
		c.JSON(http.StatusBadGateway, ExecutionResponse{
			Error:         "Failed to execute code",
			StatusMessage: "Failed to execute code",
			Success:       false,
		})
		return
	}

	var resp ExecutionResponse
	if err := json.Unmarshal(msg, &resp); err != nil {
		c.JSON(http.StatusBadRequest, ExecutionResponse{
			Error:         err.Error(),
			StatusMessage: "Failed to parse response",
			Success:       false,
		})
		return
	}

	if resp.Error == "rate_limit_exceeded" {
		var rl rateLimitResponse
		json.Unmarshal(msg, &rl)
		c.JSON(http.StatusTooManyRequests, rl)
		return
	}

	c.JSON(http.StatusOK, resp)
}
