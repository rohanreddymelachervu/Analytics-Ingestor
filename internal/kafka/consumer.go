package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/rohanreddymelachervu/ingestor/internal/models"
)

type Consumer struct {
	consumerGroup sarama.ConsumerGroup
	topics        []string
	eventService  EventProcessor
	ready         chan bool
}

// EventProcessor interface to avoid circular imports
type EventProcessor interface {
	ProcessEvent(event models.EventPayload, userID interface{}) error
}

type ConsumerGroupHandler struct {
	eventService EventProcessor
	ready        chan bool
	once         sync.Once
}

func NewConsumer(brokers []string, groupID string, topics []string, eventService EventProcessor) (*Consumer, error) {
	config := sarama.NewConfig()

	// Consumer settings for reliability
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Group.Session.Timeout = 10 * time.Second
	config.Consumer.Group.Heartbeat.Interval = 3 * time.Second
	config.Consumer.MaxProcessingTime = 1 * time.Minute

	// Enable auto-commit for simplicity
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	// Create consumer group
	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return &Consumer{
		consumerGroup: consumerGroup,
		topics:        topics,
		eventService:  eventService,
		ready:         make(chan bool),
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	handler := &ConsumerGroupHandler{
		eventService: c.eventService,
		ready:        c.ready,
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			// Check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}

			// Consume should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := c.consumerGroup.Consume(ctx, c.topics, handler); err != nil {
				log.Printf("Error from consumer: %v", err)
				return
			}
		}
	}()

	// Wait till consumer is ready
	<-c.ready
	log.Println("Kafka consumer up and running....")

	// Listen for termination signals
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		log.Println("Context cancelled, terminating consumer")
	case <-sigterm:
		log.Println("Termination signal received, terminating consumer")
	}

	wg.Wait()
	return c.consumerGroup.Close()
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready (only once)
	h.once.Do(func() {
		close(h.ready)
	})
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE: Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s, partition = %d, offset = %d",
				string(message.Value), message.Timestamp, message.Topic, message.Partition, message.Offset)

			if err := h.processMessage(message); err != nil {
				log.Printf("Error processing message: %v", err)
				// Continue processing other messages even if one fails
			}

			// Mark message as processed
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

func (h *ConsumerGroupHandler) processMessage(message *sarama.ConsumerMessage) error {
	// Parse the Kafka message
	var eventMessage EventMessage
	if err := json.Unmarshal(message.Value, &eventMessage); err != nil {
		return fmt.Errorf("failed to unmarshal event message: %w", err)
	}

	// Convert to original EventPayload format
	payloadBytes, err := json.Marshal(eventMessage.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	var eventPayload models.EventPayload
	if err := json.Unmarshal(payloadBytes, &eventPayload); err != nil {
		return fmt.Errorf("failed to unmarshal event payload: %w", err)
	}

	// Process using existing business logic
	if err := h.eventService.ProcessEvent(eventPayload, nil); err != nil {
		return fmt.Errorf("failed to process event: %w", err)
	}

	log.Printf("Successfully processed event: %s (type: %s)", eventMessage.EventID, eventMessage.EventType)
	return nil
}
