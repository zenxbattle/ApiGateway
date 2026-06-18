package clients

import (
	"context"
	"log"
	"time"
	config "xcode/configs"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type ClientConnections struct {
	ConnUser    *grpc.ClientConn
	ConnProblem *grpc.ClientConn
}

//load balancing in grpc 
//way1: we can implement app level lb using the grpc client itself that can automatically lb using any algos(default: round robin)
//way2: dedicated lb server using nginx etc that can implement much more than just being an LB - (ratelimiting)

func InitClients(config *config.Config) (*ClientConnections, error) {
	targetProblem := config.ProblemGRPCURL
	targetUser := config.UserGRPCURL

	for {
		log.Println("Connecting to ProblemGRPC:", targetProblem)
		connProblem, err := grpc.NewClient(targetProblem,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:    10 * time.Second,
				Timeout: 3 * time.Second,
			}),
		)
		if err != nil {
			log.Printf("Failed to create Problem gRPC client: %v, retrying...", err)
			time.Sleep(3 * time.Second)
			continue
		}
		connProblem.Connect()

		log.Println("Connecting to UserGRPC:", targetUser)
		connUser, err := grpc.NewClient(targetUser,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:    10 * time.Second,
				Timeout: 3 * time.Second,
			}),
		)
		if err != nil {
			connProblem.Close()
			log.Printf("Failed to create User gRPC client: %v, retrying...", err)
			time.Sleep(3 * time.Second)
			continue
		}
		connUser.Connect()

		if !waitForConnection(connProblem, 5*time.Second) || !waitForConnection(connUser, 5*time.Second) {
			connProblem.Close()
			connUser.Close()
			log.Println("gRPC connections not ready, retrying...")
			time.Sleep(3 * time.Second)
			continue
		}

		log.Println("All gRPC connections established")
		return &ClientConnections{
			ConnUser:    connUser,
			ConnProblem: connProblem,
		}, nil
	}
}

// waitForConnection checks if the gRPC connection reaches the READY state within the timeout
func waitForConnection(conn *grpc.ClientConn, timeout time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return false
		default:
			if conn.GetState() == connectivity.Ready {
				return true
			}
			time.Sleep(100 * time.Millisecond) // Poll interval
		}
	}
}

func (c *ClientConnections) Close() {
	if c.ConnUser != nil {
		c.ConnUser.Close()
	}
	if c.ConnProblem != nil {
		c.ConnProblem.Close()
	}

}
