package event

import (
	"context"
	"log"
	"time"

	"github.com/SavanRajyaguru/ecommerce-go-notification-service/config"
	"github.com/SavanRajyaguru/ecommerce-go-notification-service/internal/processor"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	processor *processor.NotificationProcessor
}

func NewConsumer(proc *processor.NotificationProcessor) *Consumer {
	return &Consumer{processor: proc}
}

func (c *Consumer) Start(ctx context.Context) {
	brokers := config.AppConfig.Kafka.Brokers
	topics := config.AppConfig.Kafka.Topics
	groupID := config.AppConfig.Kafka.GroupID

	if len(brokers) == 0 {
		log.Println("No Kafka brokers configured, skipping consumer")
		return
	}

	// We need to listen to multiple topics.
	// segmentio/kafka-go Reader can listen to a single topic or we can use Group with multiple topics?
	// ReaderConfig has "Topic" string or "GroupTopics" []string.

	topicList := []string{}
	for k, v := range topics {
		// topics map maps "order.created" -> "actual-topic-name"
		// We want to verify we are listening to what we expect.
		// For simplicity, let's collect all values from the map.
		// If map is empty, we default to standard names?
		log.Printf("Subscribing to logic topic %s: %s", k, v)
		topicList = append(topicList, v)
	}

	// If no topics mapped, add defaults
	if len(topicList) == 0 {
		topicList = []string{"order.created", "order.cancelled", "payment.success", "payment.failed"}
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupTopics: topicList,
		GroupID:     groupID,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		MaxWait:     1 * time.Second,
	})

	log.Println("Kafka Consumer started, topics:", topicList)

	for {
		// 1. Fetch Message
		m, err := r.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				break // Context cancelled
			}
			log.Printf("Error fetching message: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// 2. Process Message
		log.Printf("Received message from topic %s", m.Topic)

		// We probably need to map actual topic name back to logical name if logic depends on it
		// Or logic just depends on payload or we trust m.Topic matches our switch case in processor.

		if err := c.processor.ProcessEvent(ctx, m.Topic, m.Value); err != nil {
			log.Printf("Error processing message: %v", err)
			// Retry?
			// If transient error, maybe don't commit?
			// For now, we log and commit to avoid blocking forever on bad message.
		}

		// 3. Commit Offset
		if err := r.CommitMessages(ctx, m); err != nil {
			log.Printf("Error committing message: %v", err)
		}
	}

	r.Close()
}
