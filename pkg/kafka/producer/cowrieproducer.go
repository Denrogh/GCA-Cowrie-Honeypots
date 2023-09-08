package kafka

import (
        "encoding/json"
        "fmt"
	"time"

        "github.com/Shopify/sarama"
)

//TODO: Make into a 'model' mirror
type CowrieSession struct {
	HoneypotName     string
        SessionID        string    `json:"session"`
        StartTime        time.Time `json:"startTime"`
        EndTime          time.Time `json:"endTime"`
        PeerIP           string    `json:"peerIP"`
        PeerPort         int       `json:"peerPort"`
        HostIP           string    `json:"hostIP"`
        HostPort         int       `json:"hostPort"`
        LoggedIn         *[]string  `json:"loggedin"`
        Credentials      *[][]string  `json:"credentials"`
        Commands         *[]string  `json:"commands"`
        UnknownCommands  *[]string  `json:"unknownCommands"`
        URLs             *[]string  `json:"urls"`
        Version          *string    `json:"version"`
        TTYLog           *string   `json:"ttylog"`
        Hashes           *[]string  `json:"hashes"`
        Protocol         string    `json:"protocol"`

}

type Producer struct {
        // Kafka producer configuration
        Producer sarama.SyncProducer
}

// KafkaConfig is a struct to hold Kafka configuration.
type KafkaConfig struct {
        Brokers []string `json:"brokers"`
}

func NewProducer(kafkaBroker string) (*Producer, error) {
        // Initialize Kafka producer configuration
        config := sarama.NewConfig()
        config.Producer.Return.Successes = true

        kafkaCfg := KafkaConfig{
                Brokers: []string{kafkaBroker},
        }
        producer, err := sarama.NewSyncProducer(kafkaCfg.Brokers, config)
        if err != nil {
                return nil, fmt.Errorf("failed to create Kafka producer: %v", err)
        }

        return &Producer{
                Producer: producer,
        }, nil
}

func (p *Producer) PublishSession(session CowrieSession, topic string) error {

        // Convert CowrieSession object to JSON
        data, err := json.Marshal(session)
        if err != nil {
                return fmt.Errorf("failed to marshal session data: %v", err)

        }

	message := &sarama.ProducerMessage{
                Topic: topic,
                Value: sarama.ByteEncoder(data),
        }

        _, _, err = p.Producer.SendMessage(message)
        if err != nil {
                return fmt.Errorf("failed to produce message to Kafka: %v", err)
        }

        return nil
}

func (p *Producer) Close() {
        // Close Kafka producer
        p.Producer.Close()
}
