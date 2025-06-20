package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type Producer struct {
	syncProducer sarama.SyncProducer
	topicName    string
}

type EventMessage struct {
	EventID   string      `json:"event_id"`
	EventType string      `json:"event_type"`
	SessionID string      `json:"session_id"`
	Timestamp time.Time   `json:"timestamp"`
	Payload   interface{} `json:"payload"`
}

func NewProducer(brokers []string, topicName string) (*Producer, error) {
	config := sarama.NewConfig()

	// Producer settings for reliability
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all replicas
	config.Producer.Retry.Max = 5                    // Retry failed sends
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	// Use session_id as partition key for ordering
	config.Producer.Partitioner = sarama.NewHashPartitioner

	// Compression for better throughput
	config.Producer.Compression = sarama.CompressionSnappy

	// Batching for better performance
	config.Producer.Flush.Frequency = 100 * time.Millisecond
	config.Producer.Flush.Messages = 100

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &Producer{
		syncProducer: producer,
		topicName:    topicName,
	}, nil
}

func (p *Producer) PublishEvent(eventID, eventType, sessionID string, payload interface{}) error {
	message := EventMessage{
		EventID:   eventID,
		EventType: eventType,
		SessionID: sessionID,
		Timestamp: time.Now(),
		Payload:   payload,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	kafkaMessage := &sarama.ProducerMessage{
		Topic: p.topicName,
		Key:   sarama.StringEncoder(sessionID), // Partition by session_id
		Value: sarama.ByteEncoder(messageBytes),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("event_type"),
				Value: []byte(eventType),
			},
		},
	}

	partition, offset, err := p.syncProducer.SendMessage(kafkaMessage)
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}

	log.Printf("Event published successfully - Topic: %s, Partition: %d, Offset: %d, EventID: %s",
		p.topicName, partition, offset, eventID)

	return nil
}

func (p *Producer) Close() error {
	if p.syncProducer != nil {
		return p.syncProducer.Close()
	}
	return nil
}
