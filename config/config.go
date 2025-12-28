package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/SavanRajyaguru/ecommerce-go-config-service/proto"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	MongoURI     string
	MongoDBName  string
	SMTP         SMTPConfig
	Kafka        KafkaConfig
	FeatureFlags map[string]bool
	AppPort      string
}

type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
}

type KafkaConfig struct {
	Brokers []string          `json:"brokers"`
	Topics  map[string]string `json:"topics"` // "order.created" -> "order-created-v1", etc.
	GroupID string
}

var AppConfig *Config

func LoadConfig() {
	// 1. Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	AppConfig = &Config{
		// MongoURI is fetched from Config Service
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			Port:     getEnvInt("SMTP_PORT", 587),
			User:     getEnv("SMTP_USER", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", "no-reply@ecommerce.com"),
		},
		Kafka: KafkaConfig{
			GroupID: getEnv("KAFKA_GROUP_ID", "notification-service-group"),
		},
		AppPort: getEnv("PORT", "8086"), // Default port
	}

	if AppConfig.SMTP.Password == "" {
		log.Println("WARNING: SMTP Password not set")
	}

	// 2. Fetch from Config Service
	configServiceURL := getEnv("CONFIG_SERVICE_URL", "127.0.0.1:50051")
	log.Printf("Connecting to Config Service at: %s", configServiceURL)
	fetchRemoteConfig(configServiceURL)
}

func fetchRemoteConfig(url string) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer conn.Close()

	client := pb.NewConfigServiceClient(conn)

	// Retry logic
	var resp *pb.GetConfigResponse

	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		resp, err = client.GetConfig(ctx, &pb.GetConfigRequest{
			ServiceName: "notification-service",
		})
		cancel()

		if err == nil {
			break
		}

		log.Printf("Failed to fetch config (attempt %d/10): %v. Retrying in 2s...", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to fetch config from %s after retries: %v", url, err)
	}

	var remoteCfg struct {
		Mongo struct {
			URI      string `json:"uri"`
			Database string `json:"database"`
		} `json:"mongo"`
		Kafka struct {
			Brokers []string          `json:"brokers"`
			Topics  map[string]string `json:"topics"`
		} `json:"kafka"`
		FeatureFlags map[string]bool `json:"feature_flags"`
	}

	if err := json.Unmarshal([]byte(resp.ConfigJson), &remoteCfg); err != nil {
		log.Fatalf("Failed to unmarshal config json: %v", err)
	}

	// Mongo Config from Remote
	AppConfig.MongoURI = remoteCfg.Mongo.URI
	AppConfig.MongoDBName = remoteCfg.Mongo.Database

	// Kafka brokers -> Env variable overrides remote if present
	brokersEnv := getEnv("KAFKA_BROKERS", "")
	if brokersEnv != "" {
		AppConfig.Kafka.Brokers = []string{brokersEnv}
	} else {
		AppConfig.Kafka.Brokers = remoteCfg.Kafka.Brokers
	}

	AppConfig.Kafka.Topics = remoteCfg.Kafka.Topics
	AppConfig.FeatureFlags = remoteCfg.FeatureFlags

	fmt.Println("Configuration loaded successfully")
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		var i int
		fmt.Sscanf(value, "%d", &i)
		return i
	}
	return fallback
}
