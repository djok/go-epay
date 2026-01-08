package epay

import (
	"context"
	"time"
)

// PaymentOrderRecord represents a payment order in storage.
type PaymentOrderRecord struct {
	TransactionID string
	SubscriberID  string
	CustomerName  string
	ClientID      string
	Amount        string
	CreatedAt     time.Time
	ProcessedOn   time.Time
	InvoiceIDs    []string
}

// PaymentOrderStore is the interface for storing payment orders.
type PaymentOrderStore interface {
	// Put saves a payment order. TransactionID is used as the key.
	Put(ctx context.Context, po *PaymentOrderRecord) error

	// Get retrieves a payment order by TransactionID.
	Get(ctx context.Context, transactionID string) (*PaymentOrderRecord, error)
}
