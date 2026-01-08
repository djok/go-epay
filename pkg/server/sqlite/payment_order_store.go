package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/clouway/go-epay/pkg/epay"
	_ "github.com/mattn/go-sqlite3"
)

// PaymentOrderStore implements epay.PaymentOrderStore using SQLite.
type PaymentOrderStore struct {
	db *sql.DB
}

// NewPaymentOrderStore creates a new SQLite-backed payment order store.
func NewPaymentOrderStore(dbPath string) (*PaymentOrderStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS payment_orders (
			transaction_id TEXT PRIMARY KEY,
			subscriber_id TEXT NOT NULL,
			customer_name TEXT NOT NULL,
			client_id TEXT NOT NULL,
			amount TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			processed_on DATETIME,
			invoice_ids TEXT
		)
	`)
	if err != nil {
		db.Close()
		return nil, err
	}

	return &PaymentOrderStore{db: db}, nil
}

// Put saves a payment order.
func (s *PaymentOrderStore) Put(ctx context.Context, po *epay.PaymentOrderRecord) error {
	invoiceIDsJSON, err := json.Marshal(po.InvoiceIDs)
	if err != nil {
		return err
	}

	var processedOn *time.Time
	if !po.ProcessedOn.IsZero() {
		processedOn = &po.ProcessedOn
	}

	_, err = s.db.ExecContext(ctx, `
		INSERT OR REPLACE INTO payment_orders
		(transaction_id, subscriber_id, customer_name, client_id, amount, created_at, processed_on, invoice_ids)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, po.TransactionID, po.SubscriberID, po.CustomerName, po.ClientID, po.Amount, po.CreatedAt, processedOn, string(invoiceIDsJSON))

	return err
}

// Get retrieves a payment order by TransactionID.
func (s *PaymentOrderStore) Get(ctx context.Context, transactionID string) (*epay.PaymentOrderRecord, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT transaction_id, subscriber_id, customer_name, client_id, amount, created_at, processed_on, invoice_ids
		FROM payment_orders WHERE transaction_id = ?
	`, transactionID)

	var po epay.PaymentOrderRecord
	var invoiceIDsJSON string
	var processedOn sql.NullTime

	err := row.Scan(&po.TransactionID, &po.SubscriberID, &po.CustomerName, &po.ClientID, &po.Amount, &po.CreatedAt, &processedOn, &invoiceIDsJSON)
	if err == sql.ErrNoRows {
		return nil, epay.ErrPaymentOrderNotFound
	}
	if err != nil {
		return nil, err
	}

	if processedOn.Valid {
		po.ProcessedOn = processedOn.Time
	}

	if invoiceIDsJSON != "" {
		json.Unmarshal([]byte(invoiceIDsJSON), &po.InvoiceIDs)
	}

	return &po, nil
}

// Close closes the database connection.
func (s *PaymentOrderStore) Close() error {
	return s.db.Close()
}
