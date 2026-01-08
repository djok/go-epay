package db

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/clouway/go-epay/pkg/epay"
)

const poKind = "PaymentOrder"

// PaymentOrderStore implements epay.PaymentOrderStore using Google Cloud Datastore.
type PaymentOrderStore struct {
	client *datastore.Client
}

// NewPaymentOrderStore creates a new Datastore-backed payment order store.
func NewPaymentOrderStore(client *datastore.Client) epay.PaymentOrderStore {
	return &PaymentOrderStore{client: client}
}

type paymentOrderEntity struct {
	SubscriberID  string    `datastore:"subscriberId,noindex"`
	CustomerName  string    `datastore:"customerName,noindex"`
	ClientID      string    `datastore:"clientID,noindex"`
	TransactionID string    `datastore:"transactionId,noindex"`
	Amount        string    `datastore:"amount,noindex"`
	CreatedAt     time.Time `datastore:"createdOn,noindex"`
	ProcessedOn   time.Time `datastore:"processedOn,omitempty"`
	InvoiceIDs    []string  `datastore:"invoiceIds,noindex"`
}

// Put saves a payment order.
func (s *PaymentOrderStore) Put(ctx context.Context, po *epay.PaymentOrderRecord) error {
	k := datastore.NameKey(poKind, po.TransactionID, nil)

	entity := &paymentOrderEntity{
		SubscriberID:  po.SubscriberID,
		CustomerName:  po.CustomerName,
		ClientID:      po.ClientID,
		TransactionID: po.TransactionID,
		Amount:        po.Amount,
		CreatedAt:     po.CreatedAt,
		ProcessedOn:   po.ProcessedOn,
		InvoiceIDs:    po.InvoiceIDs,
	}

	_, err := s.client.Put(ctx, k, entity)
	return err
}

// Get retrieves a payment order by TransactionID.
func (s *PaymentOrderStore) Get(ctx context.Context, transactionID string) (*epay.PaymentOrderRecord, error) {
	k := datastore.NameKey(poKind, transactionID, nil)

	entity := &paymentOrderEntity{}
	if err := s.client.Get(ctx, k, entity); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, epay.ErrPaymentOrderNotFound
		}
		return nil, err
	}

	return &epay.PaymentOrderRecord{
		TransactionID: entity.TransactionID,
		SubscriberID:  entity.SubscriberID,
		CustomerName:  entity.CustomerName,
		ClientID:      entity.ClientID,
		Amount:        entity.Amount,
		CreatedAt:     entity.CreatedAt,
		ProcessedOn:   entity.ProcessedOn,
		InvoiceIDs:    entity.InvoiceIDs,
	}, nil
}
