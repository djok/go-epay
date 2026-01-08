package sqlite

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/clouway/go-epay/pkg/epay"
)

func TestPaymentOrderStore_PutAndGet(t *testing.T) {
	// Create temp database
	tmpFile, err := os.CreateTemp("", "test_payment_orders_*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Create store
	store, err := NewPaymentOrderStore(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	ctx := context.Background()

	// Test Put
	po := &epay.PaymentOrderRecord{
		TransactionID: "TXN123",
		SubscriberID:  "SUB456",
		CustomerName:  "Test Customer",
		ClientID:      "CLIENT789",
		Amount:        "100.50",
		CreatedAt:     time.Now(),
		InvoiceIDs:    []string{"INV1", "INV2"},
	}

	err = store.Put(ctx, po)
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	// Test Get
	retrieved, err := store.Get(ctx, "TXN123")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if retrieved.TransactionID != po.TransactionID {
		t.Errorf("TransactionID mismatch: got %s, want %s", retrieved.TransactionID, po.TransactionID)
	}
	if retrieved.SubscriberID != po.SubscriberID {
		t.Errorf("SubscriberID mismatch: got %s, want %s", retrieved.SubscriberID, po.SubscriberID)
	}
	if retrieved.CustomerName != po.CustomerName {
		t.Errorf("CustomerName mismatch: got %s, want %s", retrieved.CustomerName, po.CustomerName)
	}
	if retrieved.ClientID != po.ClientID {
		t.Errorf("ClientID mismatch: got %s, want %s", retrieved.ClientID, po.ClientID)
	}
	if retrieved.Amount != po.Amount {
		t.Errorf("Amount mismatch: got %s, want %s", retrieved.Amount, po.Amount)
	}
	if len(retrieved.InvoiceIDs) != 2 {
		t.Errorf("InvoiceIDs length mismatch: got %d, want 2", len(retrieved.InvoiceIDs))
	}

	t.Logf("Successfully stored and retrieved PaymentOrder: %+v", retrieved)
}

func TestPaymentOrderStore_GetNotFound(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_payment_orders_*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	store, err := NewPaymentOrderStore(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	ctx := context.Background()

	_, err = store.Get(ctx, "NONEXISTENT")
	if err != epay.ErrPaymentOrderNotFound {
		t.Errorf("Expected ErrPaymentOrderNotFound, got: %v", err)
	}
}

func TestPaymentOrderStore_Update(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_payment_orders_*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	store, err := NewPaymentOrderStore(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	ctx := context.Background()

	// Create initial order
	po := &epay.PaymentOrderRecord{
		TransactionID: "TXN123",
		SubscriberID:  "SUB456",
		CustomerName:  "Test Customer",
		ClientID:      "CLIENT789",
		Amount:        "100.50",
		CreatedAt:     time.Now(),
	}

	err = store.Put(ctx, po)
	if err != nil {
		t.Fatalf("Initial Put failed: %v", err)
	}

	// Update with ProcessedOn
	po.ProcessedOn = time.Now()
	err = store.Put(ctx, po)
	if err != nil {
		t.Fatalf("Update Put failed: %v", err)
	}

	// Verify update
	retrieved, err := store.Get(ctx, "TXN123")
	if err != nil {
		t.Fatalf("Get after update failed: %v", err)
	}

	if retrieved.ProcessedOn.IsZero() {
		t.Error("ProcessedOn should not be zero after update")
	}

	t.Logf("Successfully updated PaymentOrder with ProcessedOn: %v", retrieved.ProcessedOn)
}
