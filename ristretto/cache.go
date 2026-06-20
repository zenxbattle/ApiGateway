package cache

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/golang-jwt/jwt/v4"
)

type Cache struct {
	cache *ristretto.Cache
}

func NewCache() (*Cache, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,
		MaxCost:     1 << 30,
		BufferItems: 64,
	})
	if err != nil {
		return nil, err
	}
	return &Cache{cache: cache}, nil
}

// Claims structure
type Claims struct {
	ID   string `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// InvalidateToken marks a JWT token as invalidated
func (c *Cache) InvalidateToken(jwtToken string) error {
	fmt.Println("jwtToken", jwtToken)
	claims := &Claims{}
	parsedToken, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtToken), nil
	})
	if err != nil && !strings.Contains(err.Error(), "signature is invalid") {
		return fmt.Errorf("failed to parse JWT: %v", err)
	}

	var ttl time.Duration
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		expiry, ok := claims["exp"].(float64)
		if !ok {
			ttl = 24 * time.Hour
		} else {
			expiryTime := time.Unix(int64(expiry), 0)
			ttl = time.Until(expiryTime)
			if ttl <= 0 {
				ttl = time.Minute
			}
		}
	} else {
		ttl = 24 * time.Hour
	}

	ok := c.cache.SetWithTTL(jwtToken, struct{}{}, 0, ttl)
	if !ok {
		return fmt.Errorf("failed to invalidate token in cache")
	}
	// log.Printf("Cache: invalidated token: %v with TTL %v", parsedToken, ttl)
	return nil
}

// IsTokenInvalid checks if a JWT token is invalidated
func (c *Cache) IsTokenInvalid(jwtToken string) bool {
	_, found := c.cache.Get(jwtToken)
	return found
}
