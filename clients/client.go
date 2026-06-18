package clients

import (
	"context"
	"fmt"
	"log"
	"time"
	config "xcode/configs"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientConnections struct {
	ConnUser    *grpc.ClientConn
	ConnProblem *grpc.ClientConn
}

//load balancing in grpc 
//way1: we can implement app level lb using the grpc client itself that can automatically lb using any algos(default: round robin)
//way2: dedicated lb server using nginx etc that can implement much more than just being an LB - (ratelimiting)

func InitClients(config *config.Config) (*ClientConnections, error) {
	// Connect to Problem gRPC service
	targetProblem := config.ProblemGRPCURL
	log.Println("Target ProblemGRPC URL ", targetProblem)
	connProblem, err := grpc.Dial(targetProblem, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		// connUser.Close() // Clean up previous connection
		return nil, fmt.Errorf("failed to connect to Problem gRPC server: %v", err)
	}

	// Check connection state
	if !waitForConnection(connProblem, 5*time.Second) {
		// connUser.Close()
		connProblem.Close()
		return nil, fmt.Errorf("problem gRPC connection is not ready")
	}
	fmt.Println("Successfully connected to ProblemService at:", targetProblem)

	// Connect to User gRPC service
	targetUser := config.UserGRPCURL
	log.Println("Target UserGRPC URL ", targetUser)
	connUser, err := grpc.Dial(targetUser, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to User gRPC server: %v", err)
	}

	// Check connection state
	if !waitForConnection(connUser, 5*time.Second) {
		connUser.Close()
		return nil, fmt.Errorf("user gRPC connection is not ready")
	}
	fmt.Println("Successfully connected to UserService at:", targetUser)

	return &ClientConnections{
		ConnUser:    connUser,
		ConnProblem: connProblem,
	}, nil
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
