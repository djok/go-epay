package env

import (
	"context"
	"os"

	"github.com/clouway/go-epay/pkg/epay"
)

// NewEnvironmentStore creates an environment store that loads configuration
// from environment variables. This is used for Docker deployments.
func NewEnvironmentStore() epay.EnvironmentStore {
	return &envStore{}
}

type envStore struct{}

func (s *envStore) Get(ctx context.Context, name string) (*epay.Environment, error) {
	// Single-tenant mode: ignore the name parameter
	metadata := make(map[string]string)

	// UCRM-specific metadata
	if url := os.Getenv("UCRM_BILLING_URL"); url != "" {
		metadata["billingUrl"] = url
	}
	if key := os.Getenv("UCRM_API_KEY"); key != "" {
		metadata["apiKey"] = key
	}
	if id := os.Getenv("UCRM_METHOD_ID"); id != "" {
		metadata["methodId"] = id
	}
	if name := os.Getenv("UCRM_PROVIDER_NAME"); name != "" {
		metadata["providerName"] = name
	}
	if id := os.Getenv("UCRM_PROVIDER_PAYMENT_ID"); id != "" {
		metadata["providerPaymentId"] = id
	}
	if time := os.Getenv("UCRM_PROVIDER_PAYMENT_TIME"); time != "" {
		metadata["providerPaymentTime"] = time
	}
	if id := os.Getenv("UCRM_ORGANIZATION_ID"); id != "" {
		metadata["organizationId"] = id
	}

	return &epay.Environment{
		BillingJWTKey: os.Getenv("TELCONG_JWT_KEY"),
		BillingKey:    os.Getenv("TELCONG_JWT_KEY"),
		BillingURL:    os.Getenv("TELCONG_BILLING_URL"),
		EpaySecret:    os.Getenv("EPAY_SECRET"),
		MerchantID:    os.Getenv("EPAY_MERCHANT_ID"),
		Metadata:      metadata,
	}, nil
}
