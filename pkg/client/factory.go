package client

import (
	"context"
	"net/url"

	"github.com/clouway/go-epay/pkg/client/telcong"
	"github.com/clouway/go-epay/pkg/client/ucrm"
	"github.com/clouway/go-epay/pkg/epay"
	"golang.org/x/oauth2/google"
)

// BillingSystem represents the billing system type
type BillingSystem string

const (
	// BillingSystemTelcoNG represents the TelcoNG billing system
	BillingSystemTelcoNG BillingSystem = "telcong"
	// BillingSystemUCRM represents the UCRM billing system
	BillingSystemUCRM BillingSystem = "ucrm"
)

// NewClientFactory creates a new Factory for Client creation.
// This is the legacy constructor for GAE deployments that uses automatic
// billing system detection based on IDN format.
func NewClientFactory(poStore epay.PaymentOrderStore) epay.ClientFactory {
	return &clientFactory{poStore: poStore}
}

// NewClientFactoryWithBillingSystem creates a new Factory that uses the specified
// billing system for all requests. This is used for Docker deployments where the
// billing system is configured via environment variables.
func NewClientFactoryWithBillingSystem(poStore epay.PaymentOrderStore, billingSystem BillingSystem) epay.ClientFactory {
	return &clientFactory{
		poStore:       poStore,
		billingSystem: billingSystem,
	}
}

type clientFactory struct {
	poStore       epay.PaymentOrderStore
	billingSystem BillingSystem
}

func (c *clientFactory) Create(ctx context.Context, env epay.Environment, idn string) epay.Client {
	// Determine which billing system to use
	useTelcoNG := c.shouldUseTelcoNG(env, idn)

	if useTelcoNG {
		return c.createTelcoNGClient(ctx, env)
	}
	return c.createUCRMClient(env)
}

// shouldUseTelcoNG determines if TelcoNG should be used based on configuration
func (c *clientFactory) shouldUseTelcoNG(env epay.Environment, idn string) bool {
	// If billing system is explicitly set (Docker mode), use it
	if c.billingSystem != "" {
		return c.billingSystem == BillingSystemTelcoNG
	}

	// Legacy GAE mode: auto-detect based on IDN format and available config
	if IsTelcoNGContractCode(idn) && env.BillingJWTKey != "" && env.BillingURL != "" {
		return true
	}

	// If UCRM metadata is configured, use UCRM
	if _, ok := env.Metadata["billingUrl"]; ok {
		return false
	}

	// Default to TelcoNG
	return true
}

// createTelcoNGClient creates a TelcoNG billing client
func (c *clientFactory) createTelcoNGClient(ctx context.Context, env epay.Environment) epay.Client {
	billingURL, _ := url.Parse(env.BillingURL)
	conf, _ := google.JWTConfigFromJSON([]byte(env.BillingJWTKey))
	oauth2client := conf.Client(ctx)
	return telcong.NewClient(oauth2client, billingURL)
}

// createUCRMClient creates a UCRM billing client
func (c *clientFactory) createUCRMClient(env epay.Environment) epay.Client {
	billingURLStr := env.Metadata["billingUrl"]
	billingURL, _ := url.Parse(billingURLStr)
	apiKey := env.Metadata["apiKey"]
	methodID := env.Metadata["methodId"]
	providerName := env.Metadata["providerName"]
	providerPaymentID := env.Metadata["providerPaymentId"]
	providerPaymentTime := env.Metadata["providerPaymentTime"]
	organizationID := env.Metadata["organizationId"]

	return ucrm.NewClient(billingURL, apiKey, c.poStore, ucrm.PaymentProvider{
		MethodID:       methodID,
		Name:           providerName,
		PaymentID:      providerPaymentID,
		PaymentTime:    providerPaymentTime,
		OrganizationID: organizationID,
	})
}

// IsTelcoNGContractCode validates the provided code using the checksum algorithm.
// A valid code must be 7 digits with a valid Luhn checksum as the last digit.
// Delegates to epay.IsContractCode.
func IsTelcoNGContractCode(code string) bool {
	return epay.IsContractCode(code)
}
