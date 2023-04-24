package _interface

// Message represents a Pub/Sub message with specific metadata fields.
type Message struct {
	// Value contains the actual message payload (data) being sent through the Pub/Sub system.
	Value []byte

	// MessageID is a unique identifier for the message. It helps in tracking and
	// processing individual messages within the system, ensuring uniqueness and
	// preventing duplicate processing.
	MessageID string

	// Timestamp is a Unix time in milliseconds representing when the message was
	// created. This field can be useful for ordering messages, tracking message
	// age, or handling message expiration scenarios.
	Timestamp int64

	// ContentType is a string describing the content type of the message payload.
	// Examples: "application/json", "text/plain". It allows subscribers to
	// understand the format of the message and process it accordingly.
	ContentType string

	// CorrelationID is an optional field that can be used to track related
	// messages, e.g., in a request-response pattern, or to group messages
	// belonging to the same processing unit.
	CorrelationID string
}

// AckMessage represents a message acknowledgement.
type AckMessage struct {
	// MessageID is the unique identifier of the acknowledged message.
	MessageID string

	// Error contains an error value if there was an issue processing the message.
	// A nil error value indicates successful processing.
	Error error
}

// DeliveryGuarantee specifies the level of delivery guarantee for a subscription.
type DeliveryGuarantee int

const (
	// DeliverAtLeastOnce specifies that the subscriber requires at least one
	// delivery of each message. This is the most common level of guarantee.
	DeliverAtLeastOnce DeliveryGuarantee = iota

	// DeliverExactlyOnce specifies that the subscriber requires exactly one
	// delivery of each message. This guarantee level helps to prevent duplicate
	// processing in systems where processing a message multiple times could have
	// undesirable side effects.
	DeliverExactlyOnce

	// DeliverAtMostOnceWithRetry specifies that the subscriber requires at most one
	// delivery of each message, with automatic retries in case of failures.
	// This guarantee level is useful in scenarios where message delivery is
	// prioritized over ensuring that a message is processed only once.
	DeliverAtMostOnceWithRetry
)

// SubscriptionOptions contains options for a subscription.
type SubscriptionOptions struct {
	// MaxMessageSize is the maximum size of a message that the subscriber can
	// handle. It helps to prevent the subscriber from receiving messages that are
	// too large to process.
	MaxMessageSize int

	// AckTimeout is the maximum time in Unix milliseconds that the subscriber will
	// wait for an acknowledgement from the client after delivering a message.
	// If an acknowledgement is not received within this time, the message may be
	// considered as unprocessed and redelivered based on the delivery guarantee.
	AckTimeout int64

	// DeliveryGuarantee specifies the level of delivery guarantee that the
	// subscriber requires. It allows subscribers to choose the appropriate
	// level of guarantee based on their specific use case.
	DeliveryGuarantee DeliveryGuarantee

	// MaxRetries is the maximum number of retries for a failed message. It
	// ensures that the system does not get stuck in an infinite retry loop for a
	// single failed message.
	MaxRetries int
}

// Subscriber is the interface for a Pub/Sub subscriber.
type Subscriber interface {
	// Subscribe adds a handler for the specified topic. The handler function is
	// called with // a slice of messages whenever new messages are published to the topic.
	// The metadata parameter is used to provide additional information about the
	// subscription, such as the subscriber's address or connection ID.
	// The options parameter is used to specify additional subscription options,
	// such as message size limit, ack timeout, or delivery guarantee.
	// The ackChan is a channel used for sending acknowledgements for processed messages.
	Subscribe(topic string, handler func([]*Message), metadata interface{}, options *SubscriptionOptions, ackChan chan<- *AckMessage) error
	// Unsubscribe removes the handler for the specified topic. This stops the
	// subscriber from receiving messages for that topic.
	Unsubscribe(topic string) error
	// Close closes the subscriber and releases any resources. It should be called
	// when the subscriber is no longer needed to ensure proper cleanup.
	Close() error
}

// Publisher is the interface for a Pub/Sub publisher.
type Publisher interface {
	// Publish publishes the specified messages to all subscribers of the specified
	// topic. The key parameter can be used to determine the partition to which the
	// messages are sent. If the partition does not exist, it should be created.
	// The metadata parameter is used to provide additional information about the
	// messages, such as the message ID or timestamp.
	Publish(topic string, key []byte, messages []*Message, metadata interface{}) error
	// Close closes the publisher and releases any resources. It should be called
	// when the publisher is no longer needed to ensure proper cleanup.
	Close() error
}

// Partitioner is the interface for a message partitioner.
type Partitioner interface {
	// Partition returns the partition to which the specified key should be sent.
	// If the partition does not exist, it should be created. Partitions help in
	// distributing messages across different nodes or storage systems for better
	// scalability and fault tolerance.
	Partition(topic string, key []byte) (Partition, error)
}

// Partition is the interface for a message partition.
type Partition interface {
	// Subscribe adds a handler for this partition. The handler function is called
	// with a slice of messages, ordered by the ordering key. If the ordering key
	// is nil, the messages are delivered in the order in which they were published.
	Subscribe(handler func([]*Message), metadata interface{}, options *SubscriptionOptions, ackChan chan<- *AckMessage) error
	// Unsubscribe removes the handler for this partition. This stops the
	// subscriber from receiving messages for that partition.
	Unsubscribe() error

	// Publish publishes the specified message to all subscribers of this partition.
	// If the ordering key is non-nil, it should be used to determine the order in
	// which the messages are delivered to the subscribers.
	Publish(value []byte, metadata interface{}, orderingKey []byte) error

	// SetOrderingKey sets the ordering key for this partition. It is used to
	// determine the order in which messages are delivered to the subscribers.
	SetOrderingKey(orderingKey []byte) error

	// Close closes the partition and releases any resources. It should be called
	// when the partition is no longer needed to ensure proper cleanup.
	Close() error
}

// PubSub is the interface for a partitioned Pub/Sub system.
type PubSub interface {
	// Partitioner returns the partitioner for this Pub/Sub system. The partitioner
	// is responsible for determining the appropriate partition for each message,
	// ensuring better scalability and fault tolerance.
	Partitioner() Partitioner

	// Close closes the Pub/Sub system and releases any resources. It should be called
	// when the Pub/Sub system is no longer needed to ensure proper cleanup.
	Close() error
}
