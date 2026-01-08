package main

import (
	"context"
	"net/http"
	"os"

	"github.com/andyfusniak/stackdriver-gae-logrus-plugin"
	lmiddleware "github.com/andyfusniak/stackdriver-gae-logrus-plugin/middleware"
	"github.com/clouway/go-epay/pkg/client"
	"github.com/clouway/go-epay/pkg/epay"
	"github.com/clouway/go-epay/pkg/server/api"
	"github.com/clouway/go-epay/pkg/server/db"
	"github.com/clouway/go-epay/pkg/server/env"
	"github.com/clouway/go-epay/pkg/server/middleware"
	"github.com/clouway/go-epay/pkg/server/sqlite"

	"cloud.google.com/go/datastore"
	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	var envStore epay.EnvironmentStore
	var cf epay.ClientFactory

	if projectID != "" {
		// GAE Mode: Use Datastore for configuration
		log.Info("Running in GAE mode with Datastore configuration")

		formatter := stackdriver.GAEStandardFormatter(
			stackdriver.WithProjectID(projectID),
		)
		log.SetFormatter(formatter)

		dClient, err := datastore.NewClient(ctx, projectID)
		if err != nil {
			log.Fatalf("Failed to create datastore client: %v", err)
		}

		envStore = db.NewEnvironmentStore(dClient)
		poStore := db.NewPaymentOrderStore(dClient)
		cf = client.NewClientFactory(poStore)
	} else {
		// Docker Mode: Use environment variables for configuration
		log.Info("Running in Docker mode with environment variable configuration")

		// JSON formatter for Docker logs
		log.SetFormatter(&log.JSONFormatter{})

		envStore = env.NewEnvironmentStore()

		// Get billing system from environment
		billingSystem := client.BillingSystem(os.Getenv("BILLING_SYSTEM"))
		if billingSystem == "" {
			billingSystem = client.BillingSystemTelcoNG // Default to TelcoNG
		}

		// Create SQLite store for PaymentOrders (used by UCRM)
		dbPath := os.Getenv("SQLITE_DB_PATH")
		if dbPath == "" {
			dbPath = "/app/data/payment_orders.db"
		}

		var poStore epay.PaymentOrderStore
		if billingSystem == client.BillingSystemUCRM {
			var err error
			poStore, err = sqlite.NewPaymentOrderStore(dbPath)
			if err != nil {
				log.Fatalf("Failed to create SQLite store: %v", err)
			}
			log.Infof("Using SQLite database at: %s", dbPath)
		}

		cf = client.NewClientFactoryWithBillingSystem(poStore, billingSystem)
		log.Infof("Billing system configured: %s", billingSystem)
	}

	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	r := mux.NewRouter()

	// Health check endpoint for Docker
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	skipChecks := middleware.Skip("IDN", map[string]interface{}{
		"1111111111": true,
	})
	epayAPI := middleware.EpayAPIMiddleware(envStore)

	r.Handle("/v1/pay/init", skipChecks(epayAPI(api.CheckBill(cf)))).Queries("TYPE", "CHECK")
	r.Handle("/v1/pay/init", epayAPI(api.CreatePaymentOrder(cf))).Queries("TYPE", "BILLING")
	r.Handle("/v1/pay/confirm", epayAPI(api.ConfirmPaymentOrder(cf))).Queries("TYPE", "BILLING")

	http.Handle("/", lmiddleware.XCloudTraceContext(r))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
