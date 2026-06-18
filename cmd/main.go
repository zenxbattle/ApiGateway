package main

import (
	"net/http"
	"xcode/clients"
	"xcode/configs"
	metric "xcode/prometheus"
	ristretto "xcode/ristretto"
	route "xcode/route"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

func main() {
	//load configs
	cfg := configs.LoadConfig()

	//init zap logger
	appLogger, _ := zap.NewDevelopment()
	defer appLogger.Sync()

	//init gRPC clients
	grpcClient, err := clients.InitClients(&cfg)
	if err != nil {
		appLogger.Fatal("init gRPC clients failed", zap.Error(err))
	}
	defer grpcClient.Close()

	//gin router
	appRouter := gin.Default()

	//setup metrics
	requestCounter, latencyHistogram := metric.NewPrometheusClient()
	setupMetrics(appRouter, requestCounter, latencyHistogram)

	//setup middleware
	appCache, err := ristretto.NewCache()
	if err != nil {
		appLogger.Fatal("cache init failed", zap.Error(err))
	}
	requestLimiter := rate.NewLimiter(5, 15) //5 per sec refill, 15 burst
	setupMiddleware(appRouter, cfg, appLogger, appCache, requestLimiter)

	//routes
	route.SetupRoutes(appRouter, grpcClient, cfg.JWTSecretKey, appLogger)

	//start server
	port := cfg.APIGATEWAYPORT
	appLogger.Info("API Gateway running", zap.String("port", port))
	if err := appRouter.Run(":" + port); err != nil {
		appLogger.Fatal("server failed", zap.Error(err))
	}
}

// setupMetrics adds Prometheus middleware and /metrics endpoint
func setupMetrics(router *gin.Engine, requestCounter *prometheus.CounterVec, latencyHistogram *prometheus.HistogramVec) {
	router.Use(metric.PrometheusMiddleware(requestCounter, latencyHistogram))
	router.GET("/metrics", metric.Handler())
}

// setupMiddleware registers all other middleware
func setupMiddleware(router *gin.Engine, cfg configs.Config, appLogger *zap.Logger, appCache *ristretto.Cache, limiter *rate.Limiter) {
	//gin logging middleware
	router.Use(gin.Logger())

	//rate limiter
	router.Use(rateLimiterMiddleware(limiter))

	//cache middleware
	router.Use(cacheMiddleware(appCache))

	//cors setup
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-CSRF-Token", "Access-Control-Allow-Origin"},
	}))
}

// rateLimiterMiddleware returns a Gin middleware enforcing token bucket limits
func rateLimiterMiddleware(limiter *rate.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// cacheMiddleware returns a Gin middleware that sets the cache instance in context
func cacheMiddleware(appCache *ristretto.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("cacheInstance", appCache)
		c.Next()
	}
}
