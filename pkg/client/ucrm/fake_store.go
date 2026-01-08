package ucrm

import (
	"context"

	"github.com/clouway/go-epay/pkg/epay"
)

// FakePaymentOrderStore is an in-memory implementation of PaymentOrderStore for testing.
type FakePaymentOrderStore struct {
	m map[string]*epay.PaymentOrderRecord
}

// NewFakePaymentOrderStore creates a new fake store.
func NewFakePaymentOrderStore() epay.PaymentOrderStore {
	return &FakePaymentOrderStore{
		m: make(map[string]*epay.PaymentOrderRecord),
	}
}

func (s *FakePaymentOrderStore) Put(ctx context.Context, po *epay.PaymentOrderRecord) error {
	s.m[po.TransactionID] = po
	return nil
}

func (s *FakePaymentOrderStore) Get(ctx context.Context, transactionID string) (*epay.PaymentOrderRecord, error) {
	po, ok := s.m[transactionID]
	if !ok {
		return nil, epay.ErrPaymentOrderNotFound
	}
	return po, nil
}
