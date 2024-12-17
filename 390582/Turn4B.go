package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
)

type RateLimitClaim struct {
	jwt.StandardClaims
	Role      string `json:"role"`
	RateLimit int    `json:"rate_limit"`
}

var (
	secretKey   = []byte("your_secret_key")
	etcdClient  *clientv3.Client
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	configKey      = "/config/rate_limits"
	roleRateLimits = make(map[string]int)
)

// Function to initialize Etcd client
func initEtcdClient() error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://127.0.0.1:2379"}, // Replace with your Etcd endpoint
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error creating Etcd client: %v", err)
	}
	etcdClient = cli
	return nil
}

// Function to fetch rate limits from Etcd
func fetchRateLimitsFromEtcd() error {
	resp, err := etcdClient.Get(context.Background(), configKey)
	if err != nil {
		return fmt.Errorf("error fetching config from Etcd: %v", err)
	}

	if len(resp.Kvs) == 0 {
		return fmt.Errorf("config key not found in Etcd")
	}

	if err := json.Unmarshal(resp.Kvs[0].Value, &roleRateLimits); err != nil {
		return fmt.Errorf("error unmarshaling config: %v", err)
	}

	return nil
}

// Subscribe to Etcd key changes to handle real-time updates
func subscribeToEtcdChanges() {
	rch := etcdClient.Watch(context.Background(), configKey)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			if ev.Type == clientv3.EventTypePut {
				// Clear the existing rate limits
				roleRateLimits = make(map[string]int)
				if err := json.Unmarshal(ev.Kv.Value, &roleRateLimits); err != nil {
					log.Printf("error unmarshaling config update: %v", err)
				}
				log.Println("Rate limit config updated, applied new limits")
			}
		}
	}
}

// (middleware function remains the same as before)

func main() {
	if err := initEtcdClient(); err != nil {
		log.Fatalf("error initializing Etcd: %v", err)
	}
	defer etcdClient.Close()

	// Fetch initial rate limits from Etcd
	if err := fetchRateLimitsFromEtcd(); err != nil {
		log.Fatalf("error fetching initial rate limits: %v", err)
	}
}
